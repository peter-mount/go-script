package stdlib

import (
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/script"
)

func _newArray(e executor.Executor, _ *script.CallFunc) error {
	e.Calculator().Push([]interface{}{})
	return nil
}

func _append(slice []any, elems ...any) []any {
	return append(slice, elems...)
}
