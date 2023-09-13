package stdlib

import (
	"context"
	"fmt"
	error2 "github.com/peter-mount/go-script/error"
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/script"
	"reflect"
)

func _len(e executor.Executor, call *script.CallFunc, ctx context.Context) (err error) {
	var arg []interface{}
	arg, err = executor.Args(e, call, ctx)
	if err != nil {
		return err
	}
	if len(arg) != 1 {
		return fmt.Errorf("expected 1 arg, got %d", len(arg))
	}

	// Len can panic if it's not a valid type
	// so convert that to a normal error
	defer func() {
		if err1 := recover(); err1 != nil {
			err = error2.Errorf(call.Pos, "%v", err1)
		}
	}()

	tv := reflect.ValueOf(arg[0])
	ti := reflect.Indirect(tv)
	return error2.NewReturn(ti.Len())
}
