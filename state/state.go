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
	GetFunction(n string) (*script.FuncDec, bool)
	GetFunctions() []string
}

type state struct {
	mutex     sync.Mutex
	script    *script.Script
	functions map[string]*script.FuncDec
}

func New(s *script.Script) (State, error) {
	state := &state{
		script:    s,
		functions: make(map[string]*script.FuncDec),
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
