package stdlib

import (
	"context"
	"fmt"
	error2 "github.com/peter-mount/go-script/error"
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/script"
)

func _print(e executor.Executor, call *script.CallFunc, ctx context.Context) error {
	a, err := executor.Args(e, call, ctx)
	if err != nil {
		return error2.Error(call.Pos, err)
	}
	fmt.Print(a...)
	return nil
}

func _println(e executor.Executor, call *script.CallFunc, ctx context.Context) error {
	a, err := executor.Args(e, call, ctx)
	if err != nil {
		return error2.Error(call.Pos, err)
	}
	fmt.Println(a...)
	return nil
}
