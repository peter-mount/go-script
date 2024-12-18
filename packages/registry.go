package packages

import (
	"fmt"
	"reflect"
	"sync"
)

var (
	mutex    sync.Mutex
	packages = map[string]any{}
)

// RegisterPackage registers a package using its go package name as the name within g-script.
//
// Limitations:
//
// Only one package can be registered per go package. If a second package is registered for the same go package this
// will panic.
//
// Packages registered with RegisterPackage must be imported manually within a script using the import statement just
// like in go.
func RegisterPackage(f any) {
	t := reflect.TypeOf(f)
	t = t.Elem()
	Register(t.PkgPath(), f)
}

// Register a package against a name.
//
// If these packages are a single word (contains no '.' or '/' in them) then those packages are global and are
// accessible in scripts without requiring an import statement.
//
// This allows for a globally defined object similar to a package of functions or constants.
//
// This will panic if name has already been registered.
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
