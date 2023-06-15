package state

import (
	"context"
	"github.com/peter-mount/go-script/packages"
	"github.com/peter-mount/go-script/script"
	"sort"
	"strings"
	"sync"
)

const (
	stateKey = "go-script/state"
)

// State holds the current processing state of the Script
type State interface {
	Variables
	// GetFunction by name
	GetFunction(n string) (*script.FuncDec, bool)
	// GetFunctions returns a list of declared functions
	GetFunctions() []string
	WithContext(context.Context) context.Context
}

type state struct {
	mutex     sync.Mutex
	script    *script.Script
	functions map[string]*script.FuncDec
	variables Variables
}

func New(s *script.Script) (State, error) {
	state := &state{
		script:    s,
		functions: make(map[string]*script.FuncDec),
		variables: NewVariables(),
	}
	return state, state.setup()
}

func FromContext(ctx context.Context) State {
	return ctx.Value(stateKey).(*state)
}

func (s *state) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, stateKey, s)
}

func (s *state) GetFunction(n string) (*script.FuncDec, bool) {
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
	v, exists := s.variables.Get(n)
	if exists {
		return v, true
	}

	return packages.Lookup(n)
}

func (s *state) Declare(n string) {
	s.variables.Declare(n)
}

func (s *state) Set(n string, v interface{}) bool {
	return s.variables.Set(n, v)
}
