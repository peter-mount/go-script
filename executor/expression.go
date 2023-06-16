package executor

import (
	"context"
	"github.com/peter-mount/go-script/calculator"
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

	if op.Op == "=" {
		// left hand side must resolve to Ident
		var name string
		// This is messy TODO clean up this mess
		if eq := op.Left; eq != nil {
			if comp := eq.Left; comp != nil {
				if add := comp.Left; add != nil {
					if mul := add.Left; mul != nil {
						if unary := mul.Left; unary != nil {
							if unary.Right != nil {
								name = unary.Right.Ident
							}
						}
					}
				}
			}
		}
		if name == "" {
			return Errorf(op.Pos, "Assignment without target")
		}

		// Process RHS to get value
		err := e.visitor.VisitEquality(op.Right)
		if err != nil {
			return Error(op.Pos, err)
		}

		v, err := e.calculator.Peek()
		if err != nil {
			return Error(op.Pos, err)
		}

		// Implicit declare, e.g. `:=` used
		if op.Declare {
			e.state.Declare(name)
		}

		// Set the variable
		if !e.state.Set(name, v) {
			// Not set then declare it in this scope
			e.state.Declare(name)
			_ = e.state.Set(name, v)
		}

		return nil
	} else {
		return Error(op.Pos, e.visitor.VisitEquality(op.Left))
	}
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

	case op.CallFunc != nil:
		return Error(op.Pos, e.visitor.VisitCallFunc(op.CallFunc))

	case op.Null, op.Nil:
		e.calculator.Push(nil)

	case op.Ident != "":
		v, exists := e.state.Get(op.Ident)
		if !exists {
			return Errorf(op.Pos, "%q undefined", op.Ident)
		}

		// Handle arrays
		v, err := e.resolveArray(op, v)
		if err != nil {
			return Error(op.Pos, err)
		}

		// Resolve references
		if op.Pointer != nil {
			return Error(op.Pointer.Pos, e.getReference(op.Pointer, v, ctx))
		}

		// Just push variable onto stack
		e.calculator.Push(v)

	case op.SubExpression != nil:
		return Error(op.Pos, e.visitor.VisitExpression(op.SubExpression))

	case op.KeyValue != nil:
		if err := Error(op.Pos, e.visitor.VisitExpression(op.KeyValue.Value)); err != nil {
			return err
		}

		if val, err := e.calculator.Pop(); err != nil {
			return Error(op.KeyValue.Pos, err)
		} else {
			e.calculator.Push(calculator.NewKeyValue(op.KeyValue.Key, val))
		}

	}

	return nil
}
