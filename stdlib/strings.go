package stdlib

import (
	"context"
	"fmt"
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/script"
	"reflect"
)

func init() {
	executor.Register("len", _len)
}

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
			err = executor.Errorf(call.Pos, "%v", err1)
		}
	}()

	tv := reflect.ValueOf(arg[0])
	ti := reflect.Indirect(tv)
	return executor.NewReturn(ti.Len())
}
