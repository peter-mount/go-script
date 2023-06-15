package stdlib

import "github.com/peter-mount/go-script/executor"

func init() {
	executor.Register("fprint", _fprint)
	executor.Register("fprintln", _fprintln)
	executor.Register("len", _len)
	executor.Register("map", _map)
	executor.Register("print", _print)
	executor.Register("println", _println)
	executor.Register("throw", _throw)
}
