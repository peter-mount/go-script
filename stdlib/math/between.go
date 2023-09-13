package math

import (
	"context"
	"github.com/peter-mount/go-script/calculator"
	error2 "github.com/peter-mount/go-script/error"
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/script"
)

func init() {
	executor.Register("between", _between)
}

func _between(e executor.Executor, call *script.CallFunc, ctx context.Context) error {
	arg, err := executor.Args(e, call, ctx)
	if err != nil {
		return error2.Error(call.Pos, err)
	}
	if len(arg) != 3 {
		return error2.Errorf(call.Pos, "between(v,a,b)")
	}

	val, min, max := arg[0], arg[1], arg[2]

	calc := e.Calculator()
	result, ok, err := calc.Calculate(func(_ context.Context) error {
		return calc.Process(
			calculator.Push(val),
			calculator.Dup,
			calculator.Push(min),
			calculator.Op2(">="),
			calculator.Swap,
			calculator.Push(max),
			calculator.Op2("<="),
			calculator.Op2("&&"),
		)
	}, ctx)

	if err != nil {
		return error2.Error(call.Pos, err)
	}
	if ok {
		return error2.NewReturn(result)
	}

	// No result so return false
	return error2.NewReturn(false)
}
