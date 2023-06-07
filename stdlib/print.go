package stdlib

import (
	"context"
	"fmt"
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/script"
)

func init() {
	executor.Register("print", _print)
	executor.Register("println", _println)
}

func _print(e executor.Executor, call *script.CallFunc, ctx context.Context) error {
	a, err := executor.Args(e, call, ctx)
	if err != nil {
		return executor.Error(call.Pos, err)
	}
	fmt.Print(a...)
	return nil
}

func _println(e executor.Executor, call *script.CallFunc, ctx context.Context) error {
	a, err := executor.Args(e, call, ctx)
	if err != nil {
		return executor.Error(call.Pos, err)
	}
	fmt.Println(a...)
	return nil
}
