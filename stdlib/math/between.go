package math

import (
	"github.com/peter-mount/go-script/calculator"
	"github.com/peter-mount/go-script/errors"
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/script"
)

func init() {
	executor.Register("between", _between)
}

func _between(e executor.Executor, call *script.CallFunc) error {
	arg, err := executor.Args(e, call)
	if err != nil {
		return errors.Error(call.Pos, err)
	}
	if len(arg) != 3 {
		return errors.Errorf(call.Pos, "between(v,a,b)")
	}

	val, min, max := arg[0], arg[1], arg[2]

	calc := e.Calculator()
	result, ok, err := calc.Calculate(func() error {
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
	})

	if err != nil {
		return errors.Error(call.Pos, err)
	}
	if ok {
		return errors.NewReturn(result)
	}

	// No result so return false
	return errors.NewReturn(false)
}
