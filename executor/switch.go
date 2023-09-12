package executor

import (
	"context"
	"github.com/peter-mount/go-script/calculator"
	"github.com/peter-mount/go-script/script"
)

func (e *executor) switchStatement(ctx context.Context) (err error) {
	op := script.SwitchFromContext(ctx)

	// if present calculate the first expression which we will compare against the case's
	var left interface{}
	hasLeft := op.Expression != nil
	if hasLeft {
		left, err = e.calculator.MustCalculate(e.expression, op.Expression.WithContext(ctx))
		if err != nil {
			return
		}
	}

	for _, c := range op.Case {
		for _, expr := range c.Expression {
			var right interface{}
			b := false

			switch {
			case expr.String != nil:
				right = *expr.String

			case expr.Expression != nil:
				right, err = e.calculator.MustCalculate(e.expression, expr.Expression.WithContext(ctx))
				if err != nil {
					return
				}
			}

			b, err = e.switchCase(op, c, hasLeft, left, right)
			if b || err != nil {
				return err
			}
		}
	}

	// Default clause if we get to this point
	if op.Default != nil {
		return Error(op.Pos, e.visitor.VisitStatement(op.Default))
	}
	return nil
}

func (e *executor) switchCase(op *script.Switch, c *script.SwitchCase, hasLeft bool, left, right interface{}) (bool, error) {
	// Ignore errors here, they will be treated as false
	b := false
	if hasLeft {
		b, _ = calculator.Equals(left, right)
	} else {
		b, _ = calculator.GetBool(right)
	}

	if b {
		return true, Error(c.Pos, e.visitor.VisitStatement(c.Statement))
	}

	return false, nil
}
