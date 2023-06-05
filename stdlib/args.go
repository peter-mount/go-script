package stdlib

import (
	"context"
	"fmt"
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/script"
)

func Args(e executor.Executor, call *script.CallFunc, ctx context.Context) ([]interface{}, error) {
	v := e.Visitor()
	calc := e.Calculator()

	var a []interface{}
	for _, arg := range call.Args {
		v, ok, err := calc.Calculate(func(ctx context.Context) error {
			return v.VisitExpression(arg)
		}, ctx)
		switch {
		case err != nil:
			return nil, err
		case ok:
			a = append(a, v)
		default:
			return nil, fmt.Errorf("%s no value returned", arg.Pos)
		}
	}
	return a, nil
}
