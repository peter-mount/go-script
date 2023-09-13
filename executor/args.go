package executor

import (
	"context"
	"github.com/peter-mount/go-script/errors"
	"github.com/peter-mount/go-script/script"
	"math"
)

func Args(e Executor, call *script.CallFunc, ctx context.Context) ([]interface{}, error) {
	visitor := e.Visitor()
	calc := e.Calculator()

	var a []interface{}

	if call.Parameters != nil {
		for _, arg := range call.Parameters.Args {
			val, valReturned, err := calc.Calculate(func(ctx context.Context) error {
				return visitor.VisitExpression(arg)
			}, ctx)
			switch {
			case err != nil:
				return nil, errors.Error(arg.Pos, err)
			case valReturned:
				a = append(a, val)
			default:
				return nil, errors.Errorf(arg.Pos, "no value returned")
			}
		}
	}

	return a, nil
}

// RequireArgs enforces a function to be called with n arguments
func (f Function) RequireArgs(n int) Function {
	return f.RequireArgsRange(n, n)
}

// RequireMinArgs enforces a function to be called with at least n arguments
func (f Function) RequireMinArgs(n int) Function {
	return f.RequireArgsRange(n, math.MaxInt)
}

// RequireMaxArgs enforces a function to be called with a maximum of n arguments
func (f Function) RequireMaxArgs(n int) Function {
	return f.RequireArgsRange(0, n)
}

// RequireArgsRange enforces a function to be called with a range of arguments
func (f Function) RequireArgsRange(min, max int) Function {
	if min < 0 {
		min = 0
	}
	if min > max {
		min, max = max, min
	}
	return f.Then(func(e Executor, call *script.CallFunc, ctx context.Context) error {
		l := 0
		if call.Parameters != nil {
			l = len(call.Parameters.Args)
		}
		switch {
		case min == max && l != min:
			return errors.Errorf(call.Pos, "%s requires %d arguments")

		case l < min:
			return errors.Errorf(call.Pos, "%s requires minimum of %d arguments")
		case l > max:
			return errors.Errorf(call.Pos, "%s requires maximum of %d arguments")
		}
		return nil
	})
}
