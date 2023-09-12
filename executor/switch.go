package executor

import (
	"context"
	"github.com/peter-mount/go-script/calculator"
	"github.com/peter-mount/go-script/script"
)

func (e *executor) switchStatement(ctx context.Context) error {
	op := script.SwitchFromContext(ctx)

	// if present calculate the first expression which we will compare against the case's
	var left interface{}
	if op.Expression != nil {
		v, err := e.calculator.MustCalculate(e.expression, op.Expression.WithContext(ctx))
		if err != nil {
			return err
		}
		left = v
	}

	for _, c := range op.Case {
		var v interface{}
		if c.Expression != nil {
			v1, err := e.calculator.MustCalculate(e.expression, c.Expression.WithContext(ctx))
			if err != nil {
				return err
			}
			v = v1
		} else if c.String != nil {
			v = *c.String
		}

		// Ignore errors here, they will be treated as false
		b := false
		if op.Expression == nil {
			b, _ = calculator.GetBool(v)
		} else {
			b, _ = calculator.Equals(left, v)
		}

		if b {
			// return here as we don't want to have follow through
			return Error(c.Pos, e.visitor.VisitStatement(c.Statement))
		}
	}

	// Default clause if we get to this point
	if op.Default != nil {
		return Error(op.Pos, e.visitor.VisitStatement(op.Default))
	}
	return nil
}
