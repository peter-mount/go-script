package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/peter-mount/go-script/errors"
	"github.com/peter-mount/go-script/script"
)

type initialiser struct {
	state initState
}

// initState holds various state during the init scan
type initState struct {
	inLoop bool // true when parsing within a loop statement
}

func (p *defaultParser) init(s *script.Script) (*script.Script, error) {
	err := p.includeTopDec(s, s)
	if err != nil {
		return nil, err
	}

	init := &initialiser{}

	err = init.scan(s)
	if err != nil {
		return nil, errors.Error(s.Pos, err)
	}

	return s, nil
}

func (p *initialiser) scan(s *script.Script) error {
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

	return errors.Error(op.Pos, p.initStatements(op.FunBody))
}

// initStatement sets Statements.Parent to point to the Statements
// instance containing this one.
//
// # It works by presuming that the immediate caller is usually Statement.Block
//
// The only time this should not be the case is when processing a
// function definition. As that's a top level object, there is no Statement
// in the context, hence Statements.Parent will be nil.
func (p *initialiser) initStatements(op *script.Statements) error {
	if op == nil {
		return nil
	}

	for i, s := range op.Statements {
		if i > 0 {
			op.Statements[i-1].Next = s
		}

		if err := p.initStatement(s); err != nil {
			return errors.Error(s.Pos, err)
		}
	}

	return nil
}

// initStatement sets Statement.Parent to the containing Statements instance.
//
// This works as Statement is only ever contained within a Statements instance
// so the parent is its parent.
func (p *initialiser) initStatement(op *script.Statement) error {
	if op == nil {
		return nil
	}

	var err error

	switch {
	case op.Block != nil:
		err = p.initStatements(op.Block)

	case op.IfStmt != nil:
		err = p.initIf(op.IfStmt)

	case op.For != nil:
		err = p.initFor(op.For)

	case op.ForRange != nil:
		err = p.initForRange(op.ForRange)

	case op.Repeat != nil:
		err = p.initRepeat(op.Repeat)

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
	err := p.initStatement(op.Body)

	if err == nil {
		err = p.initStatement(op.Else)
	}

	return errors.Error(op.Pos, err)
}

func (p *initialiser) initSwitch(op *script.Switch) error {
	var err error

	for _, c := range op.Case {
		err = p.initStatement(c.Statement)
		if err != nil {
			break
		}
	}

	if err == nil {
		err = p.initStatement(op.Default)
	}

	return errors.Error(op.Pos, err)
}

func (p *initialiser) initDoWhile(op *script.DoWhile) error {
	return p.initLoop(op.Pos, op.Body)
}

func (p *initialiser) initRepeat(op *script.Repeat) error {
	return p.initLoop(op.Pos, op.Body)
}

func (p *initialiser) initWhile(op *script.While) error {
	return p.initLoop(op.Pos, op.Body)
}

func (p *initialiser) initFor(op *script.For) error {
	return p.initLoop(op.Pos, op.Body)
}

func (p *initialiser) initForRange(op *script.ForRange) error {
	return p.initLoop(op.Pos, op.Body)
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

	err := p.initStatement(body)

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

	err := p.initStatement(op.Body)

	if err == nil && op.Catch != nil {
		err = p.initStatement(op.Catch.Statement)
	}

	if err == nil && op.Finally != nil {
		err = p.initStatement(op.Finally.Statement)
	}

	return errors.Error(op.Pos, err)
}
