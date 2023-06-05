package executor

import (
	"context"
	"fmt"
	"github.com/peter-mount/go-script/script"
)

func (e *executor) expression(ctx context.Context) error {
	op := script.ExpressionFromContext(ctx)

	if op.Assignment != nil {
		v, exists, err := e.calculator.Calculate(e.assignment, op.Assignment.WithContext(ctx))
		if err != nil {
			return err
		}
		if exists {
			e.calculator.Push(v)
		}
	}

	return nil
}

func (e *executor) assignment(ctx context.Context) error {
	op := script.AssignmentFromContext(ctx)

	if op.Right != nil {
		// Assignment
	}

	return e.visitor.VisitEquality(op.Left)
}

func (e *executor) equality(ctx context.Context) error {
	op := script.EqualityFromContext(ctx)

	if err := e.visitor.VisitComparison(op.Left); err != nil {
		return err
	}

	if op.Right != nil {
		if err := e.visitor.VisitEquality(op.Right); err != nil {
			return err
		}
		return e.calculator.Op2(op.Op)
	}

	return nil
}

func (e *executor) comparison(ctx context.Context) error {
	op := script.ComparisonFromContext(ctx)

	if err := e.visitor.VisitAddition(op.Left); err != nil {
		return err
	}

	if op.Right != nil {
		if err := e.visitor.VisitComparison(op.Right); err != nil {
			return err
		}
		return e.calculator.Op2(op.Op)
	}

	return nil
}

func (e *executor) addition(ctx context.Context) error {
	op := script.AdditionFromContext(ctx)

	if err := e.visitor.VisitMultiplication(op.Left); err != nil {
		return err
	}

	if op.Right != nil {
		if err := e.visitor.VisitAddition(op.Right); err != nil {
			return err
		}
		return e.calculator.Op2(op.Op)
	}

	return nil
}

func (e *executor) multiplication(ctx context.Context) error {
	op := script.MultiplicationFromContext(ctx)

	if err := e.visitor.VisitUnary(op.Left); err != nil {
		return err
	}

	if op.Right != nil {
		if err := e.visitor.VisitMultiplication(op.Right); err != nil {
			return err
		}
		return e.calculator.Op2(op.Op)
	}

	return nil
}

func (e *executor) unary(ctx context.Context) error {
	op := script.UnaryFromContext(ctx)

	if op.Unary != nil {
		// TODO implement
		return fmt.Errorf("%s %q not implemented", op.Pos.String(), op.Op)
	}

	if op.Primary != nil {
		return e.visitor.VisitPrimary(op.Primary)
	}

	return nil
}

func (e *executor) primary(ctx context.Context) error {
	op := script.PrimaryFromContext(ctx)

	switch {
	case op.Float != nil:
		e.calculator.Push(*op.Float)

	case op.Integer != nil:
		e.calculator.Push(*op.Integer)

	case op.String != nil:
		e.calculator.Push(*op.String)

	case op.ArrayIndex != nil:
		return fmt.Errorf("%s ArrayIndex not implemented", op.Pos.String())

	case op.CallFunc != nil:
		return e.visitor.VisitCallFunc(op.CallFunc)

	case op.Ident != "":
		v, exists := e.state.Get(op.Ident)
		if !exists {
			return fmt.Errorf("%s %q undefined", op.Pos.String(), op.Ident)
		}
		e.calculator.Push(v)

	case op.SubExpression != nil:
		return e.visitor.VisitExpression(op.SubExpression)

	}

	return nil
}
