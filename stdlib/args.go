package stdlib

import (
	"context"
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/script"
)

func Args(e executor.Executor, call *script.CallFunc, ctx context.Context) ([]interface{}, error) {
	visitor := e.Visitor()
	calc := e.Calculator()

	var a []interface{}
	for _, arg := range call.Args {
		val, valReturned, err := calc.Calculate(func(ctx context.Context) error {
			return visitor.VisitExpression(arg)
		}, ctx)
		switch {
		case err != nil:
			return nil, executor.Error(arg.Pos, err)
		case valReturned:
			a = append(a, val)
		default:
			return nil, executor.Errorf(arg.Pos, "no value returned")
		}
	}
	return a, nil
}
