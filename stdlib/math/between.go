package math

import (
	"context"
	"github.com/peter-mount/go-script/calculator"
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/script"
	"github.com/peter-mount/go-script/stdlib"
)

func init() {
	executor.Register("between", _between)
}

func _between(e executor.Executor, call *script.CallFunc, ctx context.Context) error {
	arg, err := stdlib.Args(e, call, ctx)
	if err != nil {
		return executor.Error(call.Pos, err)
	}
	if len(arg) != 3 {
		return executor.Errorf(call.Pos, "between(v,a,b)")
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
		return executor.Error(call.Pos, err)
	}
	if ok {
		return executor.NewReturn(result)
	}

	// No result so return false
	return executor.NewReturn(false)
}
