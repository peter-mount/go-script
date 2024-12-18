package state

import (
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/peter-mount/go-script/packages"
	"github.com/peter-mount/go-script/script"
	"sort"
	"strings"
	"sync"
)

// State holds the current processing state of the Script
type State interface {
	Variables

	// GetFunction by name
	GetFunction(pos lexer.Position, n string) (*script.FuncDec, bool)

	// GetFunctions returns a list of declared functions
	GetFunctions() []string

	// SetFunction sets the current FuncDec in use, returning the previous one
	SetFunction(currentFunction *script.FuncDec) *script.FuncDec
}

type state struct {
	mutex           sync.Mutex
	script          *script.Script             // The script being executed, with included scripts
	functions       map[string]*script.FuncDec // The declared functions in all scripts
	variables       Variables                  // The current variable scope
	packages        map[string]any             // Imported packages
	currentFunction *script.FuncDec            // The function currently being executed
}

func New(s *script.Script) (State, error) {
	state := &state{
		script:    s,
		functions: make(map[string]*script.FuncDec),
		variables: NewVariables(),
		packages:  make(map[string]any),
	}
	return state, state.setup()
}

func (s *state) GetFunction(pos lexer.Position, n string) (*script.FuncDec, bool) {

	// If function name starts with _ then it's local so prefix with the filename
	// so that it's not accessible from outside its own file.
	//
	// The format of the actual name stored in the map is in a format that a
	// script cannot use for a function call, protecting the local method.
	//
	// This is only used here and in declareFunction().
	if strings.HasPrefix(n, "_") {
		n = "!" + pos.Filename + "!" + n
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	e, exists := s.functions[n]
	return e, exists
}

func (s *state) getFunctions() []string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var r []string
	for k, _ := range s.functions {
		r = append(r, k)
	}
	return r
}

func (s *state) GetFunctions() []string {
	r := s.getFunctions()
	sort.SliceStable(r, func(i, j int) bool {
		return strings.ToLower(r[i]) < strings.ToLower(r[j])
	})
	return r
}

func (s *state) NewScope() Variables {
	s.variables = s.variables.NewScope()
	return s
}

func (s *state) NewRootScope() Variables {
	s.variables = s.variables.NewRootScope()
	return s
}

func (s *state) EndScope() Variables {
	s.variables = s.variables.EndScope()
	return s
}

func (s *state) GlobalScope() Variables {
	return s.variables.GlobalScope()
}

func (s *state) Get(n string) (interface{}, bool) {
	if v, exists := s.variables.Get(n); exists {
		return v, true
	}

	// Lookup locally imported packages
	localName := s.currentFunction.Pos.Filename + "!" + n
	if pkg, exists := s.packages[localName]; exists {
		return pkg, true
	}

	// Lookup global packages - e.g. this works as n here will not contain either
	// '.' or '/' as those will not form an Ident.
	//
	// This allows for legacy package registrations to work whilst allowing the core
	// packages to always be short formed.
	return packages.Lookup(n)
}

func (s *state) Declare(n string) {
	s.variables.Declare(n)
}

func (s *state) Set(n string, v interface{}) bool {
	return s.variables.Set(n, v)
}

func (s *state) SetFunction(currentFunction *script.FuncDec) *script.FuncDec {
	old := s.currentFunction
	s.currentFunction = currentFunction
	return old
}
