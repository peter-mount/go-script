package packages

import (
	"fmt"
	"sync"
)

var (
	mutex    sync.Mutex
	packages = map[string]any{}
)

// Register an object against a name.
// This allows for a globally defined object similar to a package of functions or constants.
// This will panic if name has already been registered
func Register(name string, f any) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := packages[name]; exists {
		panic(fmt.Errorf("package %q already registered", name))
	}

	packages[name] = f
}

// Lookup a registered instance by name
func Lookup(name string) (any, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	f, exists := packages[name]
	if exists {
		return f, true
	}

	return nil, false
}
