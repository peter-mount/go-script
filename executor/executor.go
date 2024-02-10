package executor

import (
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/peter-mount/go-script/calculator"
	"github.com/peter-mount/go-script/errors"
	"github.com/peter-mount/go-script/script"
	"github.com/peter-mount/go-script/state"
	"reflect"
)

type Executor interface {
	ExpressionExecutor
	Run() error
	// ProcessParameters will call each parameter in a CallFunc returning the true values
	ProcessParameters(*script.CallFunc) ([]interface{}, error)
	// ArgsToValues will take a slice of arguments and convert to reflect.Value.
	// This will handle if CallFunc.Variadic is set
	ArgsToValues(cf *script.CallFunc, tf reflect.Type, args []interface{}) ([]reflect.Value, error)
	CallReflectFuncImpl(*script.CallFunc, reflect.Value, []interface{}) (interface{}, error)
}

type ExpressionExecutor interface {
	Calculator() calculator.Calculator
	GlobalScope() state.Variables
	Expression(op *script.Expression) error
	Statement(statements *script.Statement) error
	Statements(statements *script.Statements) error
}

type executor struct {
	script     *script.Script
	state      state.State
	calculator calculator.Calculator
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

	return e, nil
}

// NewExpressionExecutor returns an Executor that can evaluate Expressions.
// All packages are available, but no functions can be defined.
func NewExpressionExecutor() ExpressionExecutor {
	//  We can ignore error here as we do not use functions
	exec, _ := New(&script.Script{})
	return exec
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

func (e *executor) Calculator() calculator.Calculator {
	return e.calculator
}

func (e *executor) GlobalScope() state.Variables {
	return e.state.GlobalScope()
}
