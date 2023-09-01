package state

type Variables interface {
	// NewScope creates a new Variables scope
	NewScope() Variables
	// EndScope closes this Variables scope and returns the previous one.
	// If this is the Global scope this returns this instance.
	EndScope() Variables
	// NewRootScope is like NewScope except that variable lookups are not
	// passed to the parent if missing from the current scope.
	// This is used to isolate variables inside functions
	NewRootScope() Variables
	// GlobalScope returns the global scope, usually the one returned by NewVariables()
	GlobalScope() Variables
	// Declare a variable.
	Declare(n string)
	// Set a variable. If the variable is declared in a parent scope this
	// will set the variable there.
	// returns false if the variable is undeclared.
	Set(n string, val interface{}) bool
	// Get returns the variable, checking parent scopes until it finds it.
	Get(string) (interface{}, bool)
}

type variables struct {
	parent      *variables // If not nil then parent scope for variable resolving
	trueParent  *variables // Actual parent when ending a scope. nil for global
	globalScope *variables // Pointer to the root global scope
	vars        map[string]interface{}
}

func NewVariables() Variables {
	return newVariables(nil, nil)
}

func newVariables(parent, trueParent *variables) Variables {
	v := &variables{
		parent:     parent,
		trueParent: trueParent,
		vars:       make(map[string]interface{}),
	}

	if v.trueParent == nil {
		v.globalScope = v
	} else {
		v.globalScope = v.trueParent.globalScope
	}

	return v
}

func (v *variables) NewScope() Variables {
	return newVariables(v, v)
}

func (v *variables) NewRootScope() Variables {
	return newVariables(v.globalScope, v)
}

func (v *variables) EndScope() Variables {
	if v.trueParent == nil {
		return v
	}
	return v.trueParent
}

func (v *variables) GlobalScope() Variables {
	return v.globalScope
}

func (v *variables) Get(n string) (interface{}, bool) {
	if r, exists := v.vars[n]; exists {
		return r, true
	}
	if v.parent != nil {
		return v.parent.Get(n)
	}
	return nil, false
}

func (v *variables) Set(n string, val interface{}) bool {
	if _, exists := v.vars[n]; exists {
		v.vars[n] = val
		return true
	}

	// If a sub-scope then check the parent
	if v.parent != nil {
		return v.parent.Set(n, val)
	}

	return false
}

func (v *variables) Declare(n string) {
	if IsValidVariable(n) {
		v.vars[n] = nil
	}
}

func IsValidVariable(n string) bool {
	return n != "" && n != "_"
}
