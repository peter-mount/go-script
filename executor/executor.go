package executor

import (
	"context"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/peter-mount/go-script/calculator"
	"github.com/peter-mount/go-script/errors"
	"github.com/peter-mount/go-script/script"
	"github.com/peter-mount/go-script/state"
	"github.com/peter-mount/go-script/visitor"
	"reflect"
)

type Executor interface {
	Run() error
	Visitor() visitor.Visitor
	Calculator() calculator.Calculator
	GlobalScope() state.Variables
	// ProcessParameters will call each parameter in a CallFunc returning the true values
	ProcessParameters(*script.CallFunc, context.Context) ([]interface{}, error)
	// ArgsToValues will take a slice of arguments and convert to reflect.Value.
	// This will handle if CallFunc.Variadic is set
	ArgsToValues(cf *script.CallFunc, tf reflect.Type, args []interface{}) ([]reflect.Value, error)
	CallReflectFuncImpl(*script.CallFunc, reflect.Value, []interface{}) (interface{}, error)
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
		CallFunc(e.callFunc).
		DoWhile(e.doWhile).
		Expression(e.expression).
		For(e.forStatement).
		ForRange(e.forRange).
		If(e.ifStatement).
		Repeat(e.repeatUntil).
		Return(e.returnStatement).
		Statement(e.statement).
		Statements(e.statements).
		StatementsNoNest().
		Switch(e.switchStatement).
		Try(e.try).
		While(e.while).
		WithContext(e.context)

	return e, nil
}

func (e *executor) Run() error {
	main, hasMain := e.state.GetFunction(lexer.Position{}, "main")
	if !hasMain {
		return errors.Errorf(e.script.Pos, "main() function not defined")
	}

	err := e.functionImpl(main, nil)

	// Pass err unless it's return, break or continue.
	// break should happen lower down but this catches it, so it doesn't
	// exit the function call
	if err != nil && !(errors.IsReturn(err) || errors.IsBreak(err)) || errors.IsContinue(err) {
		return errors.Error(e.script.Pos, err)
	}

	return nil
}

func (e *executor) Visitor() visitor.Visitor {
	return e.visitor
}

func (e *executor) Calculator() calculator.Calculator {
	return e.calculator
}

func (e *executor) GlobalScope() state.Variables {
	return e.state.GlobalScope()
}
