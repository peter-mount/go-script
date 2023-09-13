package stdlib

import (
	"context"
	"fmt"
	"github.com/peter-mount/go-script/calculator"
	error2 "github.com/peter-mount/go-script/error"
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/script"
)

func _throw(e executor.Executor, call *script.CallFunc, ctx context.Context) error {
	a, err := executor.Args(e, call, ctx)
	if err == nil {

		switch len(a) {
		case 0:
			return fmt.Errorf("throw(format[,...])")

		case 1:
			if err1, ok := a[0].(error); ok {
				// expression is an error so use it
				err = err1
			} else if message, err1 := calculator.GetString(a[0]); err1 != nil {
				err = err1
			} else {
				err = error2.Errorf(call.Pos, message)
			}

		default:
			if format, err1 := calculator.GetString(a[0]); err1 != nil {
				err = err1
			} else {
				err = error2.Errorf(call.Pos, format, a[1:]...)
			}
		}
	}

	return error2.Error(call.Pos, err)
}
