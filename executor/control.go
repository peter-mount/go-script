package executor

import (
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/peter-mount/go-script/calculator"
	"github.com/peter-mount/go-script/errors"
	"github.com/peter-mount/go-script/script"
	"github.com/peter-mount/go-script/state"
	"reflect"
)

// condition evaluates Expression, and converts the result to a boolean.
// defaultResult is the result returned if expr is nil.
func (e *executor) condition(expr *script.Expression, defaultResult bool) (bool, error) {
	if expr == nil {
		return defaultResult, nil
	}

	v, err := e.calculator.MustCalculate(func() error {
		return e.Expression(expr)
	})
	if err != nil {
		return false, err
	}

	b, err := calculator.GetBool(v)
	return b, errors.Error(expr.Pos, err)
}

// breakOrContinue checks for errors, break and continue statements.
// the bool is true if the loop should be terminated, false to continue or no error.
// error will be the error to return, or nil if no error.
func (e *executor) breakOrContinue(pos lexer.Position, err error) (bool, error) {
	// Consume break and exit the loop
	if errors.IsBreak(err) {
		return true, nil
	}

	// Consume continue
	if errors.IsContinue(err) {
		return false, nil
	}

	// a normal error
	return err != nil, errors.Error(pos, err)
}

func (e *executor) ifStatement(s *script.If) error {

	b, err := e.condition(s.Condition, true)
	if err == nil {
		if b {
			err = e.Statement(s.Body)
		} else {
			err = e.Statement(s.Else)
		}
	}

	return errors.Error(s.Pos, err)
}

// repeatUntil from basic etc. repeats body until condition is met.
// body is always evaluated once.
func (e *executor) repeatUntil(s *script.Repeat) error {
	return e.forLoop(s.Pos, nil, nil, s.Body, nil, s.Condition, false)
}

// doWhile from C, repeats body while condition is met.
// body is always executed once.
func (e *executor) doWhile(s *script.DoWhile) error {
	return e.forLoop(s.Pos, nil, nil, s.Body, nil, s.Condition, true)
}

// while from C, execute body while condition is met.
// body will never run if condition never passes
func (e *executor) while(s *script.While) error {
	return e.forLoop(s.Pos, nil, s.Condition, s.Body, nil, nil, true)
}

// forStatement from C, optional init & increment but executes body while condition is met.
// body will never run if condition never passes.
func (e *executor) forStatement(s *script.For) error {
	return e.forLoop(s.Pos, s.Init, s.Condition, s.Body, s.Increment, nil, true)
}

// forLoop is the internals of loops.
// p is the Position of the statement being implemented.
// init is the optional init Expression
// conditionFirst is the condition test performed at the start of the loop
// body the Statement to execute inside the loop
// inc is the optional increment Expression
// conditionLast is the condition test performed at the end of the loop
// conditionResult the result of conditionFirst or conditionLast to repeat the loop.
func (e *executor) forLoop(p lexer.Position, init, conditionFirst *script.Expression, body *script.Statement, inc, conditionLast *script.Expression, conditionResult bool) error {

	// Run for in a new scope so variables declared there are not accessible outside
	e.state.NewScope()
	defer e.state.EndScope()

	if init != nil {
		err := e.Expression(init)
		if err != nil {
			return errors.Error(p, err)
		}
	}

	for {
		b, err := e.condition(conditionFirst, conditionResult)
		if err != nil || b != conditionResult {
			return errors.Error(p, err)
		}

		if body != nil {
			exit, err1 := e.breakOrContinue(p, e.Statement(body))
			if exit {
				return err1
			}
		}

		if inc != nil {
			err = e.Expression(inc)
			if err != nil {
				return errors.Error(p, err)
			}
		}

		// conditionResult false then condition last
		b, err = e.condition(conditionLast, conditionResult)
		if err != nil || b != conditionResult {
			return errors.Error(p, err)
		}
	}
}

func (e *executor) forRange(op *script.ForRange) error {
	// Run for in a new scope so variables declared there are not accessible outside
	e.state.NewScope()
	defer e.state.EndScope()

	// Declare in scope if := used
	if op.Declare {
		e.state.Declare(op.Key)
		e.state.Declare(op.Value)
	}

	// Evaluate Expression
	r, err := e.calculator.MustCalculate(func() error { return e.Expression(op.Expression) })
	if err != nil {
		return errors.Error(op.Pos, err)
	}

	// Check for supported interfaces
	if r != nil {
		// If an Iterator then run through until HasNext() returns false
		if it, ok := r.(Iterator); ok {
			return e.forIterator(op, it)
		}
	}

	// Handle default go constructs
	tv := reflect.ValueOf(r)
	ti := reflect.Indirect(tv)
	switch ti.Kind() {
	case reflect.Map:
		return e.forMapIter(op, ti.MapRange())

	case reflect.Array, reflect.Slice, reflect.String:
		return e.forSlice(op, ti)

	default:
		return errors.Errorf(op.Expression.Pos, "cannot range over %T", r)
	}
}

// forIterator will iterate for all values in an Iterator
func (e *executor) forIterator(op *script.ForRange, it Iterator) error {
	for i := 0; it.HasNext(); i++ {
		exit, err := e.breakOrContinue(op.Pos, e.forRangeEntry(i, it.Next(), op))
		if exit {
			return err
		}
	}
	return nil
}

// forMapIter will iterate over a MapIter
func (e *executor) forMapIter(op *script.ForRange, mi *reflect.MapIter) error {
	for mi.Next() {
		exit, err := e.breakOrContinue(op.Pos, e.forRangeEntry(mi.Key().Interface(), mi.Value().Interface(), op))
		if exit {
			return err
		}
	}
	return nil
}

// forSlice will iterate over a reflect.Array, reflect.Slice or reflect.String.
// This will panic if Value is not one of those types.
func (e *executor) forSlice(op *script.ForRange, ti reflect.Value) error {
	l := ti.Len()
	for i := 0; i < l; i++ {
		exit, err := e.breakOrContinue(op.Pos, e.forRangeEntry(i, ti.Index(i).Interface(), op))
		if exit {
			return err
		}
	}
	return nil
}

// forRangeEntryImpl is used in for range either from forRangeEntryValue or an iterator
func (e *executor) forRangeEntry(key, val interface{}, op *script.ForRange) error {
	if op.Body == nil {
		return nil
	}

	if state.IsValidVariable(op.Key) {
		if !e.state.Set(op.Key, key) {
			e.state.Declare(op.Key)
			_ = e.state.Set(op.Key, key)
		}
	}

	if state.IsValidVariable(op.Value) {
		if !e.state.Set(op.Value, val) {
			e.state.Declare(op.Value)
			_ = e.state.Set(op.Value, val)
		}
	}

	return errors.Error(op.Pos, e.Statement(op.Body))
}

type Iterator interface {
	HasNext() bool
	Next() interface{}
}
