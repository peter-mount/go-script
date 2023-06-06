package executor

import (
	"context"
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
		return Error(s.Pos, e.visitor.VisitStatement(s.Body))
	} else {
		return Error(s.Pos, e.visitor.VisitStatement(s.Else))
	}
}

func (e *executor) whileStatement(ctx context.Context) error {
	s := script.WhileFromContext(ctx)

	for {
		b, err := e.condition(s.Condition, ctx)
		if err != nil {
			return Error(s.Pos, err)
		}

		if !b {
			return nil
		}

		if err := Error(s.Pos, e.visitor.VisitStatement(s.Body)); err != nil {
			return Error(s.Pos, err)
		}
	}
}
