package calculator

import (
	"context"
	"fmt"
	"github.com/peter-mount/go-kernel/v2/util/task"
	"strings"
)

const (
	calculatorKey = "go-script/calculatorKey"
)

type Calculator interface {
	// WithContext adds the Calculator to a Context
	WithContext(ctx context.Context) context.Context
	Reset() Calculator
	// Push a value onto the stack
	Push(v interface{}) Calculator
	// Pop a value from the stack. Return an error if the stack is empty
	// Returns an error if the stack is empty.
	Pop() (interface{}, error)
	// Pop2 returns two values from the stack.
	// The order returned (a,b) is if b was the top value on the stack whilst a was below it.
	// Returns an error if the stack is empty.
	Pop2() (interface{}, interface{}, error)
	// Peek returns the top value on the stack, the stack is unchanged
	// Returns an error if the stack is empty.
	Peek() (interface{}, error)
	// Swap the top two entries on the stack.
	// Return an error if the stack does not have two entries to swap.
	Swap() error
	// Rot rotates third itm to top [n3 n2 n1] -> [n2 n1 n3]
	Rot() error
	// Over copies second item to top [n2 n1] -> [n2 n1 n2]
	Over() error
	// Dup duplicates the top entry on the stack
	// Returns an error if the stack is empty.
	Dup() error
	// Drop removes the top entry on the stack
	// Returns an error if the stack is empty.
	Drop() error
	// Op2 performs a named operation using the top two entries on the stack, returning
	// the result to the stack.
	//
	// Returns an error if the stack doesn't have two entries or the calculation fails
	Op2(op string) error
	// Calculate executes a Calculation against this calculator.
	//
	// It ensures that the stack is valid for just this calculation, preserving
	// any existing stack for the calculation that is executing this one.
	// Returns an error if the calculation fails.
	//
	// Returns the top value of the stack at the end of the calculation,
	// or nil if the stack was empty.
	Calculate(t task.Task, ctx context.Context) (interface{}, bool, error)
	// Exec is similar to Calculate but does not return a result.
	Exec(t task.Task, ctx context.Context) error
	// Process processes a series of Instruction's to perform a calculation
	Process(instructions ...Instruction) error
	// Dump returns the current stack as a string.
	// Used for debugging
	Dump() string
}

// FromContext returns the Calculator associated with a Context
func FromContext(ctx context.Context) Calculator {
	c := ctx.Value(calculatorKey)
	if c != nil {
		return c.(*calculator)
	}
	return nil
}

// New Calculator
func New() Calculator {
	return &calculator{}
}

type calculator struct {
	stack []interface{}
}

func (c *calculator) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, calculatorKey, c)
}

func (c *calculator) Reset() Calculator {
	c.stack = nil
	return c
}

func (c *calculator) Push(v interface{}) Calculator {
	c.stack = append(c.stack, v)
	return c
}

func (c *calculator) Pop() (interface{}, error) {
	if len(c.stack) == 0 {
		return nil, stackEmpty
	}
	l := len(c.stack) - 1
	v := c.stack[l]
	c.stack = c.stack[:l]
	return v, nil
}

func (c *calculator) Pop2() (interface{}, interface{}, error) {
	b, err := c.Pop()
	if err != nil {
		return nil, nil, err
	}

	a, err := c.Pop()
	if err != nil {
		return nil, nil, err
	}

	return a, b, nil
}

func (c *calculator) Peek() (interface{}, error) {
	if len(c.stack) == 0 {
		return nil, stackEmpty
	}
	l := len(c.stack) - 1
	v := c.stack[l]
	return v, nil
}

func (c *calculator) Swap() error {
	l := len(c.stack)
	if l < 2 {
		return stackEmpty
	}
	c.stack[l-2], c.stack[l-1] = c.stack[l-1], c.stack[l-2]
	return nil
}

func (c *calculator) Dup() error {
	l := len(c.stack)
	if l < 1 {
		return stackEmpty
	}
	c.stack = append(c.stack, c.stack[l-1])
	return nil
}

// Rot [n3 n2 n1] -> [n2 n1 n3]
func (c *calculator) Rot() error {
	l := len(c.stack)
	if l < 3 {
		return stackEmpty
	}
	c.stack[l-3], c.stack[l-2], c.stack[l-1] = c.stack[l-2], c.stack[l-1], c.stack[l-3]
	return nil
}

// Over [n2 n1] -> [n2 n1 n2]
func (c *calculator) Over() error {
	l := len(c.stack)
	if l < 2 {
		return stackEmpty
	}
	c.stack = append(c.stack, c.stack[l-2])
	return nil
}

func (c *calculator) Drop() error {
	_, err := c.Pop()
	return err
}

func (c *calculator) Op2(op string) error {
	operation, exists := operations[op]
	if !exists {
		return fmt.Errorf("operation %q undefined", op)
	}

	a, b, err := c.Pop2()
	if err != nil {
		return err
	}

	v, err := operation.BiCalculate(a, b)
	if err != nil {
		return err
	}

	c.Push(v)
	return nil
}

func (c *calculator) Calculate(t task.Task, ctx context.Context) (interface{}, bool, error) {
	if c == nil {
		return nil, false, nil
	}

	old := c.stack
	c.stack = nil
	defer func() {
		c.stack = old
	}()

	err := t(ctx)
	if err != nil {
		return nil, false, err
	}

	// Ignore the error as a value is optional but use it for the return bool
	v, err := c.Pop()
	return v, err == nil, nil
}

func (c *calculator) Exec(t task.Task, ctx context.Context) error {
	if c == nil {
		return nil
	}

	old := c.stack
	c.stack = nil
	defer func() {
		c.stack = old
	}()

	return t.Do(ctx)
}

func (c *calculator) Dump() string {
	var a []string
	for _, e := range c.stack {
		a = append(a, fmt.Sprintf("%T[%v]", e, e))
	}
	return "Stack=[" + strings.Join(a, " ") + "]"
}
