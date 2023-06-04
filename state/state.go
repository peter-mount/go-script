package state

import (
	"context"
	"fmt"
	"github.com/peter-mount/go-script/script"
	"github.com/peter-mount/go-script/visitor"
	"sort"
	"strings"
	"sync"
)

// State holds the current processing state of the Script
type State interface {
	Variables
	// GetFunction by name
	GetFunction(n string) (*script.FuncDec, bool)
	// GetFunctions returns a list of declared functions
	GetFunctions() []string
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

	err := visitor.New().
		FuncDec(state.declareFunction).
		WithContext(context.Background()).
		VisitScript(s)

	return state, err
}

func (s *state) declareFunction(ctx context.Context) error {
	f := script.FuncDecFromContext(ctx)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if e, exists := s.functions[f.Name]; exists {
		return fmt.Errorf("%s function %q already defined at %s", f.Pos.String(), f.Name, e.Pos.String())
	}
	s.functions[f.Name] = f
	return nil
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
func (s *state) EndScope() Variables {
	s.variables = s.variables.EndScope()
	return s
}

func (s *state) Declare(n string) { s.variables.Declare(n) }

func (s *state) Set(n string, v interface{}) bool { return s.variables.Set(n, v) }

func (s *state) Get(n string) (interface{}, bool) { return s.variables.Get(n) }
