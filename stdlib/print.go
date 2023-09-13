package stdlib

import (
	"fmt"
	"github.com/peter-mount/go-script/errors"
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/script"
)

func _print(e executor.Executor, call *script.CallFunc) error {
	a, err := executor.Args(e, call)
	if err != nil {
		return errors.Error(call.Pos, err)
	}
	fmt.Print(a...)
	return nil
}

func _println(e executor.Executor, call *script.CallFunc) error {
	a, err := executor.Args(e, call)
	if err != nil {
		return errors.Error(call.Pos, err)
	}
	fmt.Println(a...)
	return nil
}
