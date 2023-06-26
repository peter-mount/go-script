package executor

import (
	"context"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/peter-mount/go-script/calculator"
	"github.com/peter-mount/go-script/script"
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

func (e *executor) whileStatement(ctx context.Context) error {
	s := script.WhileFromContext(ctx)

	// while cond body is the same as for ;cond; body
	return e.forLoop(s.Pos, nil, s.Condition, nil, s.Body, ctx)
}

func (e *executor) forStatement(ctx context.Context) error {
	s := script.ForFromContext(ctx)
	return e.forLoop(s.Pos, s.Init, s.Condition, s.Increment, s.Body, ctx)
}

func (e *executor) forLoop(p lexer.Position, init, condition, inc *script.Expression, body *script.Statement, ctx context.Context) error {

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
		if condition != nil {
			b, err := e.condition(condition, ctx)
			if err != nil {
				return Error(p, err)
			}

			if !b {
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

				return Error(p, err)
			}
		}

		if inc != nil {
			if err := e.visitor.VisitExpression(inc); err != nil {
				return Error(p, err)
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
		if op.Key != "_" {
			e.state.Declare(op.Key)
		}
		if op.Value != "_" {
			e.state.Declare(op.Value)
		}
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

				return Error(op.Pos, err)
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

				return Error(op.Pos, err)
			}
		}

	default:
		return Errorf(op.Expression.Pos, "cannot range over %T", r)
	}

	return nil
}

func (e *executor) forRangeEntry(key, val reflect.Value, op *script.ForRange, ctx context.Context) error {
	if op.Key != "_" {
		if !e.state.Set(op.Key, key.Interface()) {
			e.state.Declare(op.Key)
			_ = e.state.Set(op.Key, key.Interface())
		}
	}

	if op.Value != "_" {
		if !e.state.Set(op.Value, val.Interface()) {
			e.state.Declare(op.Value)
			_ = e.state.Set(op.Value, val.Interface())
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
