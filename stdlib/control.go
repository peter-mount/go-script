package stdlib

import (
	"fmt"
	"github.com/peter-mount/go-script/calculator"
	"github.com/peter-mount/go-script/errors"
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/script"
)

func _throw(e executor.Executor, call *script.CallFunc) error {
	a, err := executor.Args(e, call)
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
				err = errors.Errorf(call.Pos, message)
			}

		default:
			if format, err1 := calculator.GetString(a[0]); err1 != nil {
				err = err1
			} else {
				err = errors.Errorf(call.Pos, format, a[1:]...)
			}
		}
	}

	return errors.Error(call.Pos, err)
}
