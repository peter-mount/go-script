package executor

import (
	"context"
	"fmt"
	"github.com/peter-mount/go-script/script"
	"sync"
)

var (
	mutex   sync.Mutex
	library = map[string]Function{}
)

type Function func(e Executor, call *script.CallFunc, ctx context.Context) error

func Register(name string, f Function) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := library[name]; exists {
		panic(fmt.Errorf("function %q already registered", name))
	}

	library[name] = f
}

func Lookup(name string) (Function, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	f, exists := library[name]
	if exists {
		return f, true
	}

	return nil, false
}
