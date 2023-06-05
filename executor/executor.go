package executor

import (
	"context"
	"fmt"
	"github.com/peter-mount/go-script/calculator"
	"github.com/peter-mount/go-script/script"
	"github.com/peter-mount/go-script/state"
	"github.com/peter-mount/go-script/visitor"
)

type Executor interface {
	Run() error
}

type executor struct {
	script     *script.Script
	state      state.State
	calculator calculator.Calculator
	visitor    visitor.Visitor
}

func New(s *script.Script) (Executor, error) {
	execState, err := state.New(s)
	if err != nil {
		return nil, err
	}
	e := &executor{
		script:     s,
		state:      execState,
		calculator: calculator.New(),
	}
	e.visitor = visitor.New().
		Statement(e.statement).
		WithContext(execState.WithContext(context.Background()))

	return e, nil
}

func (e *executor) Run() error {
	main, hasMain := e.state.GetFunction("main")
	if !hasMain {
		return fmt.Errorf("%s main() function not defined", e.script.Pos)
	}

	fmt.Printf("exec %q\n", main.Name)

	return nil
}
