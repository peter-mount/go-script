package executor

import (
	"github.com/peter-mount/go-script/calculator"
	"github.com/peter-mount/go-script/errors"
	"github.com/peter-mount/go-script/script"
)

func (e *executor) switchStatement(op *script.Switch) (err error) {
	// if present calculate the first Expression which we will compare against the case's
	var left interface{}
	hasLeft := op.Expression != nil
	if hasLeft {
		left, err = e.calculator.MustCalculate(func() error { return e.Expression(op.Expression) })
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
				right, err = e.calculator.MustCalculate(func() error { return e.Expression(expr.Expression) })
				if err != nil {
					return
				}
			}

			b, err = e.switchCase(c, hasLeft, left, right)
			if b || err != nil {
				return err
			}
		}
	}

	// Default clause if we get to this point
	if op.Default != nil {
		return errors.Error(op.Pos, e.Statement(op.Default))
	}
	return nil
}

func (e *executor) switchCase(c *script.SwitchCase, hasLeft bool, left, right interface{}) (bool, error) {
	// Ignore errors here, they will be treated as false
	b := false
	if hasLeft {
		b, _ = calculator.Equals(left, right)
	} else {
		b, _ = calculator.GetBool(right)
	}

	if b {
		return true, errors.Error(c.Pos, e.Statement(c.Statement))
	}

	return false, nil
}
