package executor

import (
	"context"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/peter-mount/go-script/calculator"
	"github.com/peter-mount/go-script/script"
	"github.com/peter-mount/go-script/state"
	"reflect"
)

func (e *executor) condition(expr *script.Expression, ctx context.Context) (bool, error) {
	v, ok, err := e.calculator.Calculate(e.expression, expr.WithContext(ctx))
	if err != nil {
		return false, Error(expr.Pos, err)
	}
	if !ok {
		return false, Errorf(expr.Pos, "No result from condition")
	}

	b, err := calculator.GetBool(v)
	return b, Error(expr.Pos, err)
}

func (e *executor) ifStatement(ctx context.Context) error {
	s := script.IfFromContext(ctx)

	b, err := e.condition(s.Condition, ctx)
	if err != nil {
		return Error(s.Pos, err)
	}

	if b {
		err = Error(s.Pos, e.visitor.VisitStatement(s.Body))
	} else {
		err = Error(s.Pos, e.visitor.VisitStatement(s.Else))
	}

	return err
}

func (e *executor) repeatStatement(ctx context.Context) error {
	s := script.RepeatFromContext(ctx)

	// repeat body until cond is the same as "for ; !cond ; body"
	return e.forLoop(s.Pos, nil, nil, nil, s.Condition, s.Body, ctx, true)
}

func (e *executor) whileStatement(ctx context.Context) error {
	s := script.WhileFromContext(ctx)

	// while cond body is the same as "for ; cond ; body"
	return e.forLoop(s.Pos, nil, s.Condition, nil, nil, s.Body, ctx, true)
}

func (e *executor) forStatement(ctx context.Context) error {
	s := script.ForFromContext(ctx)
	return e.forLoop(s.Pos, s.Init, s.Condition, s.Increment, nil, s.Body, ctx, true)
}

// forLoop is the internals of loops.
// p is the Position of the statement being implemented.
// init is the optional init expression
// condition is the optional condition expression
// inc is the optional increment expression
// body the Statement to execute inside the loop
// conditionFirst true then condition tested before body executes, false then tested after body & inc
// conditionResult the result of condition to repeat the loop.
func (e *executor) forLoop(p lexer.Position, init, conditionFirst, inc, conditionLast *script.Expression, body *script.Statement, ctx context.Context, conditionResult bool) error {

	// Run for in a new scope so variables declared there are not accessible outside
	e.state.NewScope()
	defer e.state.EndScope()

	if init != nil {
		err := e.visitor.VisitExpression(init)
		if err != nil {
			return Error(p, err)
		}
	}

	for {
		// conditionResult true then condition first
		if conditionFirst != nil {
			b, err := e.condition(conditionFirst, ctx)
			if err != nil {
				return Error(p, err)
			}

			if b != conditionResult {
				return nil
			}
		}

		if body != nil {
			err := Error(p, e.visitor.VisitStatement(body))
			if err != nil {
				// Consume break and exit the loop
				if IsBreak(err) {
					return nil
				}

				// Only exit if it's not continue
				if !IsContinue(err) {
					return Error(p, err)
				}
			}
		}

		if inc != nil {
			if err := e.visitor.VisitExpression(inc); err != nil {
				return Error(p, err)
			}
		}

		// conditionResult false then condition last
		if conditionLast != nil {
			b, err := e.condition(conditionLast, ctx)
			if err != nil {
				return Error(p, err)
			}

			if b == conditionResult {
				return nil
			}
		}
	}
}

func (e *executor) forRange(ctx context.Context) error {
	op := script.ForRangeFromContext(ctx)

	// Run for in a new scope so variables declared there are not accessible outside
	e.state.NewScope()
	defer e.state.EndScope()

	// Declare in scope if := used
	if op.Declare {
		e.state.Declare(op.Key)
		e.state.Declare(op.Value)
	}

	// Evaluate expression
	err := e.visitor.VisitExpression(op.Expression)
	if err != nil {
		return Error(op.Pos, err)
	}

	r, err := e.calculator.Pop()
	if err != nil {
		return Error(op.Pos, err)
	}

	// Check for supported interfaces
	if r != nil {
		// If Iterable then convert to an Interator
		if it, ok := r.(Iterable); ok {
			r = it.Iterator()
		}

		// If an Iterator then run through until HasNext() returns false
		if it, ok := r.(Iterator); ok {
			for i := 0; it.HasNext(); i++ {
				if err := e.forRangeEntryImpl(i, it.Next(), op, ctx); err != nil {
					// Consume break and exit the loop
					if IsBreak(err) {
						return nil
					}

					// Only exit if not continue
					if !IsContinue(err) {
						return Error(op.Pos, err)
					}
				}
			}
			return nil
		}
	}

	// Handle default go constructs
	tv := reflect.ValueOf(r)
	ti := reflect.Indirect(tv)
	switch ti.Kind() {
	case reflect.Map:
		mi := ti.MapRange()
		for mi.Next() {
			if err := e.forRangeEntry(mi.Key(), mi.Value(), op, ctx); err != nil {
				// Consume break and exit the loop
				if IsBreak(err) {
					return nil
				}

				// Only exit if not continue
				if !IsContinue(err) {
					return Error(op.Pos, err)
				}
			}
		}

	case reflect.Array, reflect.Slice, reflect.String:
		l := ti.Len()
		for i := 0; i < l; i++ {
			if err := e.forRangeEntry(reflect.ValueOf(i), ti.Index(i), op, ctx); err != nil {
				// Consume break and exit the loop
				if IsBreak(err) {
					return nil
				}

				// Only exit if not continue
				if !IsContinue(err) {
					return Error(op.Pos, err)
				}
			}
		}

	default:
		return Errorf(op.Expression.Pos, "cannot range over %T", r)
	}

	return nil
}

func (e *executor) forRangeEntry(key, val reflect.Value, op *script.ForRange, ctx context.Context) error {
	return e.forRangeEntryImpl(key.Interface(), val.Interface(), op, ctx)
}

func (e *executor) forRangeEntryImpl(key, val interface{}, op *script.ForRange, _ context.Context) error {
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

	if op.Body != nil {
		err := Error(op.Pos, e.visitor.VisitStatement(op.Body))
		if err != nil {
			return Error(op.Pos, err)
		}
	}

	return nil
}

type Iterable interface {
	Iterator() Iterator
}

type Iterator interface {
	HasNext() bool
	Next() interface{}
}
