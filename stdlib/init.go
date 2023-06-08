package stdlib

import "github.com/peter-mount/go-script/executor"

func init() {
	executor.Register("len", _len)
	executor.Register("print", _print)
	executor.Register("println", _println)
	executor.Register("throw", _throw)
}
