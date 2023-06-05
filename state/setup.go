package state

import (
	"context"
	"fmt"
	"github.com/peter-mount/go-script/script"
	"github.com/peter-mount/go-script/visitor"
)

func (s *state) setup() error {
	return visitor.New().
		FuncDec(s.declareFunction).
		WithContext(context.Background()).
		VisitScript(s.script)
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
