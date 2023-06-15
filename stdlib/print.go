package stdlib

import (
	"context"
	"fmt"
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/script"
	"io"
)

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

func _fprint(e executor.Executor, call *script.CallFunc, ctx context.Context) error {
	a, err := executor.Args(e, call, ctx)
	if err != nil {
		return executor.Error(call.Pos, err)
	}

	if len(a) < 1 {
		return executor.Errorf(call.Pos, "No args")
	}

	if w, ok := a[0].(io.Writer); ok {
		_, err = fmt.Fprint(w, a[1:]...)
		return err
	}

	return executor.Errorf(call.Pos, "Expected io.Writer, got %T", a[0])
}

func _fprintln(e executor.Executor, call *script.CallFunc, ctx context.Context) error {
	a, err := executor.Args(e, call, ctx)
	if err != nil {
		return executor.Error(call.Pos, err)
	}

	if len(a) < 1 {
		return executor.Errorf(call.Pos, "No args")
	}

	if w, ok := a[0].(io.Writer); ok {
		_, err = fmt.Fprintln(w, a[1:]...)
		return err
	}

	return executor.Errorf(call.Pos, "Expected io.Writer, got %T", a[0])
}
