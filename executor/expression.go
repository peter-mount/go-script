package executor

import (
	"context"
	"github.com/peter-mount/go-script/script"
)

func (e *executor) expression(ctx context.Context) error {
	op := script.ExpressionFromContext(ctx)

	if op.Right != nil {
		v, exists, err := e.calculator.Calculate(e.assignment, op.Right.WithContext(ctx))
		if err != nil {
			return Error(op.Pos, err)
		}
		if exists {
			e.calculator.Push(v)
		}
	}

	return nil
}

func (e *executor) assignment(ctx context.Context) error {
	op := script.AssignmentFromContext(ctx)

	// Visit right hand side
	if err := e.visitor.VisitEquality(op.Right); err != nil {
		return Error(op.Pos, err)
	}

	// If Ident is set then set that variable to the result of the right hand side.
	// Peek the value, so it remains on the stack as we can use that to support multiple assignments
	// or assignments within an expression (not C but it should in theory work)
	if op.Ident != "" {
		v, err := e.calculator.Peek()
		if err != nil {
			return Error(op.Pos, err)
		}
		// Set the variable
		set := e.state.Set(op.Ident, v)
		if !set {
			// Not set then declare it in this scope
			e.state.Declare(op.Ident)
			_ = e.state.Set(op.Ident, v)
		}
	}

	return nil
}

func (e *executor) equality(ctx context.Context) error {
	op := script.EqualityFromContext(ctx)

	if err := e.visitor.VisitComparison(op.Left); err != nil {
		return Error(op.Pos, err)
	}

	if op.Right != nil {
		err := e.visitor.VisitEquality(op.Right)
		if err == nil {
			err = e.calculator.Op2(op.Op)
		}
		return Error(op.Pos, err)
	}

	return nil
}

func (e *executor) comparison(ctx context.Context) error {
	op := script.ComparisonFromContext(ctx)

	if err := e.visitor.VisitAddition(op.Left); err != nil {
		return Error(op.Pos, err)
	}

	if op.Right != nil {
		err := e.visitor.VisitComparison(op.Right)
		if err == nil {
			err = e.calculator.Op2(op.Op)
		}
		return Error(op.Pos, err)
	}

	return nil
}

func (e *executor) addition(ctx context.Context) error {
	op := script.AdditionFromContext(ctx)

	if err := e.visitor.VisitMultiplication(op.Left); err != nil {
		return Error(op.Pos, err)
	}

	if op.Right != nil {
		err := e.visitor.VisitAddition(op.Right)
		if err == nil {
			err = e.calculator.Op2(op.Op)
		}
		return Error(op.Pos, err)
	}

	return nil
}

func (e *executor) multiplication(ctx context.Context) error {
	op := script.MultiplicationFromContext(ctx)

	if err := e.visitor.VisitUnary(op.Left); err != nil {
		return Error(op.Pos, err)
	}

	if op.Right != nil {
		err := e.visitor.VisitMultiplication(op.Right)
		if err == nil {
			err = e.calculator.Op2(op.Op)
		}
		return Error(op.Pos, err)
	}

	return nil
}

func (e *executor) unary(ctx context.Context) error {
	op := script.UnaryFromContext(ctx)

	if op.Left != nil {
		// TODO implement
		return Errorf(op.Pos, "%q not implemented", op.Op)
	}

	if op.Right != nil {
		return Error(op.Pos, e.visitor.VisitPrimary(op.Right))
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
		return Errorf(op.Pos, "ArrayIndex not implemented")

	case op.CallFunc != nil:
		return Error(op.Pos, e.visitor.VisitCallFunc(op.CallFunc))

	case op.Ident != "":
		v, exists := e.state.Get(op.Ident)
		if !exists {
			return Errorf(op.Pos, "%q undefined", op.Ident)
		}
		e.calculator.Push(v)

	case op.SubExpression != nil:
		return Error(op.Pos, e.visitor.VisitExpression(op.SubExpression))

	}

	return nil
}