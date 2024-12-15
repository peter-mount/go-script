package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/peter-mount/go-script/errors"
	"github.com/peter-mount/go-script/script"
)

type Initialiser interface {
	Scan(s *script.Script) error
	Expression(op *script.Expression) error
	Statements(op *script.Statements) error
	Statement(op *script.Statement) error
}

func NewInitialiser() Initialiser {
	return &initialiser{}
}

type initialiser struct {
	state initState
}

// initState holds various state during the init Scan
type initState struct {
	inLoop bool // true when parsing within a loop statement
}

func (p *defaultParser) init(s *script.Script, err error) (*script.Script, error) {
	if err != nil {
		return nil, err
	}

	err = p.includeTopDec(s, s)
	if err != nil {
		return nil, err
	}

	init := NewInitialiser()

	err = init.Scan(s)
	if err != nil {
		return nil, errors.Error(s.Pos, err)
	}

	return s, nil
}

func (p *initialiser) Scan(s *script.Script) error {
	for _, f := range s.FunDec {
		if err := p.funcDec(f); err != nil {
			return errors.Error(s.Pos, err)
		}
	}
	return nil
}

// funcDec initialises a function declaration.
//
// Here we save the current state and set the state to a blank slate.
// The previous state is restored when we exit
func (p *initialiser) funcDec(op *script.FuncDec) error {
	old := p.state
	defer func() { p.state = old }()
	p.state = initState{}

	return errors.Error(op.Pos, p.Statements(op.FunBody))
}

// Expression initialises an expression.
func (p *initialiser) Expression(op *script.Expression) error {
	// Currently do nothing, just here to make the API usable if we ever need to include it
	return nil
}

func (p *initialiser) Statements(op *script.Statements) error {
	if op == nil {
		return nil
	}

	for i, s := range op.Statements {
		if i > 0 {
			op.Statements[i-1].Next = s
		}

		if err := p.Statement(s); err != nil {
			return errors.Error(s.Pos, err)
		}
	}

	return nil
}

func (p *initialiser) Statement(op *script.Statement) error {
	if op == nil {
		return nil
	}

	var err error

	switch {
	case op.Block != nil:
		err = p.Statements(op.Block)

	case op.Expression != nil:
		err = p.Expression(op.Expression)

	case op.DoWhile != nil:
		err = p.initDoWhile(op.DoWhile)

	case op.IfStmt != nil:
		err = p.initIf(op.IfStmt)

	case op.For != nil:
		err = p.initFor(op.For)

	case op.ForRange != nil:
		err = p.initForRange(op.ForRange)

	case op.Repeat != nil:
		err = p.initRepeat(op.Repeat)

	case op.Return != nil:
		err = p.Expression(op.Return.Result)

	case op.Switch != nil:
		err = p.initSwitch(op.Switch)

	case op.Try != nil:
		err = p.initTry(op.Try)

	case op.While != nil:
		err = p.initWhile(op.While)

	case op.Break:
		// break is only valid within a loop so this will force
		// an error if it's found outside of one
		if !p.state.inLoop {
			err = errors.Errorf(op.Pos, "break not allowed here")
		}

	case op.Continue:
		// continue is only valid within a loop so this will force
		// an error if it's found outside of one
		if !p.state.inLoop {
			err = errors.Errorf(op.Pos, "continue not allowed here")
		}
	}

	return errors.Error(op.Pos, err)
}

func (p *initialiser) initIf(op *script.If) error {
	err := p.Expression(op.Condition)

	if err == nil {
		err = p.Statement(op.Body)
	}

	if err == nil {
		err = p.Statement(op.Else)
	}

	return errors.Error(op.Pos, err)
}

func (p *initialiser) initSwitch(op *script.Switch) error {
	var err error

	for _, c := range op.Case {
		for _, ex := range c.Expression {
			err = p.Expression(ex.Expression)
			if err != nil {
				break
			}
		}
		if err == nil {
			err = p.Statement(c.Statement)
		}
		if err != nil {
			break
		}
	}

	if err == nil {
		err = p.Statement(op.Default)
	}

	return errors.Error(op.Pos, err)
}

func (p *initialiser) initDoWhile(op *script.DoWhile) error {
	err := p.Expression(op.Condition)
	if err == nil {
		err = p.initLoop(op.Pos, op.Body)
	}
	return err
}

func (p *initialiser) initRepeat(op *script.Repeat) error {
	err := p.Expression(op.Condition)
	if err == nil {
		err = p.initLoop(op.Pos, op.Body)
	}
	return err
}

func (p *initialiser) initWhile(op *script.While) error {
	err := p.Expression(op.Condition)
	if err == nil {
		err = p.initLoop(op.Pos, op.Body)
	}
	return err
}

func (p *initialiser) initFor(op *script.For) error {
	err := p.Expression(op.Init)
	if err == nil {
		err = p.Expression(op.Condition)
	}
	if err == nil {
		err = p.Expression(op.Increment)
	}
	if err == nil {
		err = p.initLoop(op.Pos, op.Body)
	}
	return err
}

func (p *initialiser) initForRange(op *script.ForRange) error {
	err := p.Expression(op.Expression)
	if err == nil {
		err = p.initLoop(op.Pos, op.Body)
	}
	return err
}

// initLoop handles all loop statements.
//
// It sets inLoop to true to mark that break/continue are now valid
// and then initialises the statement.
//
// inLoop is restored to its previous value afterwards
func (p *initialiser) initLoop(pos lexer.Position, body *script.Statement) error {
	old := p.state
	defer func() { p.state = old }()
	p.state.inLoop = true

	err := p.Statement(body)

	return errors.Error(pos, err)
}

func (p *initialiser) initTry(op *script.Try) error {

	// try-resources ensure only assignments and enforce declare mode
	// as those variables can only be accessed from within the body
	if op.Init != nil {
		for _, init := range op.Init.Resources {
			if init.Right != nil && init.Right.Right != nil {
				init.Right.Declare = true
			}
		}
	}

	err := p.Statement(op.Body)

	if err == nil && op.Catch != nil {
		err = p.Statement(op.Catch.Statement)
	}

	if err == nil && op.Finally != nil {
		err = p.Statement(op.Finally.Statement)
	}

	return errors.Error(op.Pos, err)
}
