package executor

import (
	"context"
	"github.com/peter-mount/go-script/calculator"
	"github.com/peter-mount/go-script/script"
	"github.com/peter-mount/go-script/state"
	"github.com/peter-mount/go-script/visitor"
)

type Executor interface {
	Run() error
	Visitor() visitor.Visitor
	Calculator() calculator.Calculator
}

type executor struct {
	script     *script.Script
	state      state.State
	calculator calculator.Calculator
	visitor    visitor.Visitor
	context    context.Context
}

func New(s *script.Script) (Executor, error) {
	execState, err := state.New(s)
	if err != nil {
		return nil, err
	}
	e := &executor{
		script:     s,
		state:      execState,
		calculator: calculator.New(),
	}

	e.context = execState.WithContext(context.Background())

	e.visitor = visitor.New().
		Addition(e.addition).
		Assignment(e.assignment).
		CallFunc(e.callFunc).
		Comparison(e.comparison).
		Equality(e.equality).
		Expression(e.expression).
		ExpressionNoNest().
		If(e.ifStatement).
		Multiplication(e.multiplication).
		Primary(e.primary).
		Statement(e.statement).
		Statements(e.statements).
		StatementsNoNest().
		Unary(e.unary).
		While(e.whileStatement).
		WithContext(e.context)

	return e, nil
}

func (e *executor) Run() error {
	main, hasMain := e.state.GetFunction("main")
	if !hasMain {
		return Errorf(e.script.Pos, "main() function not defined")
	}

	err := e.function(main)

	// Pass err unless it's return or break.
	// break should happen lower down but this catches it, so it doesn't
	// exit the function call
	if err != nil && !(IsReturn(err) || IsBreak(err)) {
		return Error(e.script.Pos, err)
	}

	return nil
}

func (e *executor) Visitor() visitor.Visitor {
	return e.visitor
}

func (e *executor) Calculator() calculator.Calculator {
	return e.calculator
}
