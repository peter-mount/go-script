package stdlib

import "github.com/peter-mount/go-script/executor"

func init() {
	executor.Register("append", executor.FuncDelegate(_append))
	executor.Register("len", _len)
	executor.Register("map", _map)
	executor.Register("newArray", _newArray)
	executor.Register("print", _print)
	executor.Register("println", _println)
	executor.Register("throw", _throw)
}
