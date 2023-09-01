package executor

import (
	"context"
	"github.com/peter-mount/go-script/calculator"
	"github.com/peter-mount/go-script/script"
	"reflect"
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

		var primary *script.Primary
		// This is messy TODO clean up this mess
		if log := op.Left; log != nil {
			if eq := log.Left; eq != nil {
				if comp := eq.Left; comp != nil {
					if add := comp.Left; add != nil {
						if mul := add.Left; mul != nil {
							if unary := mul.Left; unary != nil {
								primary = unary.Right
							}
						}
					}
				}
			}
		}

		if primary == nil || primary.Ident == nil || primary.Ident.Ident == "" {
			return Errorf(op.Pos, "Assignment without target")
		}

		name := primary.Ident.Ident

		if primary.Pointer == nil {
			// POVS = plain old variable setter

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
		} else {
			v, err := e.resolveIdent(primary, ctx)
			if err != nil && !IsNoFieldErr(err) {
				return Error(op.Pos, err)
			}

			if nfe, ok := GetNoFieldErr(err); ok {
				v = nfe.Value()
				name = nfe.Name()
				err = nil
			}

			vV := reflect.ValueOf(v)
			if vV.IsNil() {
				return Errorf(op.Pos, "Cannot set nil")
			}
			vT := vV.Type()

			// Process RHS to get value
			err = e.visitor.VisitEquality(op.Right)
			if err != nil {
				return Error(op.Pos, err)
			}

			setV, err := e.calculator.Peek()
			if err != nil {
				return Error(op.Pos, err)
			}

			switch vT.Kind() {
			case reflect.Map:
				vV.SetMapIndex(reflect.ValueOf(name), reflect.ValueOf(setV))

			case reflect.Struct:
				f := vV.FieldByName(name)
				if f.IsValid() && f.CanSet() {
					f.Set(reflect.ValueOf(setV))
				} else {
					return Errorf(op.Pos, "Cannot set %q on %T", name, setV)
				}

			default:
				return Errorf(op.Pos, "Cannot set %T", setV)
			}
		}

		return nil
	} else {
		return Error(op.Pos, e.visitor.VisitLogic(op.Left))
	}
}

func (e *executor) logic(ctx context.Context) error {
	op := script.LogicFromContext(ctx)

	err := e.visitor.VisitEquality(op.Left)
	for err == nil && op.Right != nil {
		err = e.visitor.VisitEquality(op.Right.Left)
		if err == nil {
			err = e.calculator.Op2(op.Op)
		}
		if err == nil {
			op = op.Right
		}
	}

	if err != nil {
		return Error(op.Pos, err)
	}
	return nil
}

func (e *executor) equality(ctx context.Context) error {
	op := script.EqualityFromContext(ctx)

	err := e.visitor.VisitComparison(op.Left)
	for err == nil && op.Right != nil {
		err = e.visitor.VisitComparison(op.Right.Left)
		if err == nil {
			err = e.calculator.Op2(op.Op)
		}
		if err == nil {
			op = op.Right
		}
	}

	if err != nil {
		return Error(op.Pos, err)
	}
	return nil
}

func (e *executor) comparison(ctx context.Context) error {
	op := script.ComparisonFromContext(ctx)

	err := e.visitor.VisitAddition(op.Left)
	for err == nil && op.Right != nil {
		err = e.visitor.VisitAddition(op.Right.Left)
		if err == nil {
			err = e.calculator.Op2(op.Op)
		}
		if err == nil {
			op = op.Right
		}
	}

	if err != nil {
		return Error(op.Pos, err)
	}
	return nil
}

func (e *executor) addition(ctx context.Context) error {
	op := script.AdditionFromContext(ctx)

	err := e.visitor.VisitMultiplication(op.Left)
	for err == nil && op.Right != nil {
		err = e.visitor.VisitMultiplication(op.Right.Left)
		if err == nil {
			err = e.calculator.Op2(op.Op)
		}
		if err == nil {
			op = op.Right
		}
	}

	if err != nil {
		return Error(op.Pos, err)
	}
	return nil
}

func (e *executor) multiplication(ctx context.Context) error {
	op := script.MultiplicationFromContext(ctx)

	err := e.visitor.VisitUnary(op.Left)
	for err == nil && op.Right != nil {
		err = e.visitor.VisitUnary(op.Right.Left)
		if err == nil {
			err = e.calculator.Op2(op.Op)
		}
		if err == nil {
			op = op.Right
		}
	}

	if err != nil {
		return Error(op.Pos, err)
	}
	return nil
}

func (e *executor) unary(ctx context.Context) error {
	op := script.UnaryFromContext(ctx)

	if op.Left != nil {
		err := e.visitor.VisitUnary(op.Left)
		if err == nil {
			err = e.calculator.Op1(op.Op)
		}
		if err != nil {
			return Error(op.Pos, err)
		}
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

	case op.True:
		e.calculator.Push(true)

	case op.False:
		e.calculator.Push(false)

	case op.Ident != nil && op.Ident.Ident != "":
		v, err := e.resolveIdent(op, ctx)
		if err != nil {
			return Error(op.Pos, err)
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

func (e *executor) resolveIdent(op *script.Primary, ctx context.Context) (interface{}, error) {
	v, exists := e.state.Get(op.Ident.Ident)
	if !exists {
		return nil, Errorf(op.Pos, "%q undefined", op.Ident)
	}

	// Handle arrays
	v, err := e.resolveArray(op, v)

	if err == nil && op.Pointer != nil {
		// Resolve references
		v, err = e.getReferenceImpl(op.Pointer, v, ctx)
	}

	if IsNoFieldErr(err) {
		return v, err
	}

	if err != nil {
		return nil, Error(op.Pos, err)
	}

	return v, nil
}
