package state

import (
	"context"
	"fmt"
	"github.com/peter-mount/go-script/script"
	"github.com/peter-mount/go-script/visitor"
	"strings"
)

func (s *state) setup() error {
	return visitor.New().
		FuncDec(s.declareFunction).
		WithContext(context.Background()).
		VisitScript(s.script)
}

func (s *state) declareFunction(ctx context.Context) error {
	f := script.FuncDecFromContext(ctx)

	// If function name starts with _ then it's local so prefix with the filename
	// so that it's not accessible from outside its own file.
	//
	// The format of the actual name stored in the map is in a format that a
	// script cannot use for a function call, protecting the local method.
	//
	// This is only used here and in GetFunction().
	name := f.Name
	if strings.HasPrefix(f.Name, "_") {
		name = "!" + f.Pos.Filename + "!" + f.Name
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if e, exists := s.functions[name]; exists {
		return fmt.Errorf("%s function %q already defined at %s", f.Pos.String(), f.Name, e.Pos.String())
	}
	s.functions[name] = f
	return nil
}
