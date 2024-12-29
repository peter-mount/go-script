package stdlib

import (
	"fmt"
	"github.com/peter-mount/go-script/calculator"
	"github.com/peter-mount/go-script/errors"
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/script"
)

// _throw implements throw(error)
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

// _isNull implements isNull builtin which returns true if any of its arguments is null/nil
func _isNull(e executor.Executor, call *script.CallFunc) error {
	a, err := executor.Args(e, call)
	if err == nil {

		switch len(a) {
		case 0:
			return fmt.Errorf("isNull(value[,...])")

		case 1:
			e.Calculator().Push(a[0] == nil)

		default:
			found := false
			for _, v := range a {
				if v == nil {
					found = true
					break
				}
			}
			e.Calculator().Push(found)
		}
	}

	return errors.Error(call.Pos, err)
}

// _notNull implements notNull builtin which returns true if all of its arguments are not null/nil
func _notNull(e executor.Executor, call *script.CallFunc) error {
	a, err := executor.Args(e, call)
	if err == nil {

		switch len(a) {
		case 0:
			return fmt.Errorf("notNull(value[,...])")

		case 1:
			e.Calculator().Push(a[0] != nil)

		default:
			// Presume all are true, fail on the first null found
			found := true
			for _, v := range a {
				if v == nil {
					found = false
					break
				}
			}
			e.Calculator().Push(found)
		}
	}

	return errors.Error(call.Pos, err)
}
