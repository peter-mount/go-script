package executor

import (
	"fmt"
	"github.com/peter-mount/go-script/calculator"
	"github.com/peter-mount/go-script/script"
	"github.com/peter-mount/go-script/state"
)

type Executor interface {
	Run() error
}

type exec struct {
	script     *script.Script
	state      state.State
	calculator calculator.Calculator
}

func New(s *script.Script) (Executor, error) {
	execState, err := state.New(s)
	if err != nil {
		return nil, err
	}
	return &exec{
		script:     s,
		state:      execState,
		calculator: calculator.New(),
	}, nil
}

func (e *exec) Run() error {
	main, hasMain := e.state.GetFunction("main")
	if !hasMain {
		return fmt.Errorf("%s main() function not defined", e.script.Pos)
	}

	fmt.Printf("exec %q\n", main.Name)

	return nil
}
