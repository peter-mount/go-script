package executor

import (
	"context"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/peter-mount/go-script/calculator"
	"github.com/peter-mount/go-script/script"
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

	// Eat break so it just exits the If
	if IsBreak(err) {
		return nil
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
