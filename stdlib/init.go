package stdlib

import "github.com/peter-mount/go-script/executor"

func init() {
	executor.Register("append", executor.FuncDelegate(_append))
	executor.Register("isNull", _isNull)
	executor.Register("len", _len)
	executor.Register("map", _map)
	executor.Register("mapContains", _mapContains)
	executor.Register("newArray", _newArray)
	executor.Register("notNull", _notNull)
	executor.Register("print", _print)
	executor.Register("println", _println)
	executor.Register("throw", _throw)
}
