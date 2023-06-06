package state

type Variables interface {
	// NewScope creates a new Variables scope
	NewScope() Variables
	// EndScope closes this Variables scope and returns the previous one.
	// If this is the Global scope this returns this instance.
	EndScope() Variables
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
	parent Variables // If not nil then parent scope
	vars   map[string]interface{}
}

func NewVariables() Variables {
	return newVariables(nil)
}

func newVariables(parent Variables) Variables {
	return &variables{
		parent: parent,
		vars:   make(map[string]interface{}),
	}
}

func (v *variables) NewScope() Variables {
	return newVariables(v)
}

func (v *variables) EndScope() Variables {
	if v.parent == nil {
		return v
	}
	return v.parent
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
	v.vars[n] = nil
}
