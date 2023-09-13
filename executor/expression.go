package executor

import (
	"github.com/peter-mount/go-script/calculator"
	"github.com/peter-mount/go-script/errors"
	"github.com/peter-mount/go-script/script"
	"reflect"
)

func (e *executor) Expression(op *script.Expression) error {

	if op.Right != nil {
		v, exists, err := e.calculator.Calculate(func() error {
			return e.assignment(op.Right)
		})
		if err != nil {
			return errors.Error(op.Pos, err)
		}
		if exists {
			e.calculator.Push(v)
		}
	}

	return nil
}

func (e *executor) assignment(op *script.Assignment) error {
	if op.Op == "=" {

		var primary *script.Primary
		// This is messy TODO clean up this mess
		if ternary := op.Left; ternary != nil {
			if log := ternary.Left; log != nil {
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
		}

		if primary == nil || primary.Ident == nil || primary.Ident.Ident == "" {
			return errors.Errorf(op.Pos, "Assignment without target")
		}

		name := primary.Ident.Ident

		if primary.Pointer == nil {
			// POVS = plain old variable setter

			// Process RHS to get value
			err := e.assignment(op.Right)
			if err != nil {
				return errors.Error(op.Pos, err)
			}

			v, err := e.calculator.Peek()
			if err != nil {
				return errors.Error(op.Pos, err)
			}

			// Augmented assignment
			if op.AugmentedOp != nil {
				v0, ok := e.state.Get(name)
				if !ok {
					return errors.Errorf(op.Pos, "%q undefined", name)
				}

				// calculate existing op v to get the true new value
				calc := e.calculator
				calc.Push(v0)
				calc.Push(v)
				err1 := calc.Op2(*op.AugmentedOp)
				if err1 == nil {
					v, err1 = calc.Pop()
				}
				if err1 != nil {
					return errors.Error(op.Pos, err1)
				}
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
			v, err := e.resolveIdent(primary)
			if err != nil && !errors.IsNoFieldErr(err) {
				return errors.Error(op.Pos, err)
			}

			if nfe, ok := errors.GetNoFieldErr(err); ok {
				v = nfe.Value()
				name = nfe.Name()
				err = nil
			}

			vV := reflect.ValueOf(v)
			if vV.IsNil() {
				return errors.Errorf(op.Pos, "Cannot set nil")
			}
			vT := vV.Type()

			// Process RHS to get value
			err = e.assignment(op.Right)
			if err != nil {
				return errors.Error(op.Pos, err)
			}

			setV, err := e.calculator.Peek()
			if err != nil {
				return errors.Error(op.Pos, err)
			}

			switch vT.Kind() {
			case reflect.Map:
				vV.SetMapIndex(reflect.ValueOf(name), reflect.ValueOf(setV))

			case reflect.Struct:
				f := vV.FieldByName(name)
				if f.IsValid() && f.CanSet() {
					f.Set(reflect.ValueOf(setV))
				} else {
					return errors.Errorf(op.Pos, "Cannot set %q on %T", name, setV)
				}

			default:
				return errors.Errorf(op.Pos, "Cannot set %T", setV)
			}
		}

		return nil
	} else {
		return errors.Error(op.Pos, e.ternary(op.Left))
	}
}

func (e *executor) ternary(op *script.Ternary) (err error) {

	err = e.logic(op.Left)
	if err == nil && op.True != nil && op.False != nil {
		var v interface{}
		v, err = e.calculator.Pop()
		if err == nil {
			var b bool
			b, err = calculator.GetBool(v)
			if err == nil {
				if b {
					err = e.logic(op.True)
				} else {
					err = e.logic(op.False)
				}
			}
		}
	}
	return errors.Error(op.Pos, err)
}

func (e *executor) logic(op *script.Logic) error {

	err := e.equality(op.Left)
	for err == nil && op.Right != nil {
		err = e.equality(op.Right.Left)
		if err == nil {
			err = e.calculator.Op2(op.Op)
		}
		if err == nil {
			op = op.Right
		}
	}

	if err != nil {
		return errors.Error(op.Pos, err)
	}
	return nil
}

func (e *executor) equality(op *script.Equality) error {

	err := e.comparison(op.Left)
	for err == nil && op.Right != nil {
		err = e.comparison(op.Right.Left)
		if err == nil {
			err = e.calculator.Op2(op.Op)
		}
		if err == nil {
			op = op.Right
		}
	}

	if err != nil {
		return errors.Error(op.Pos, err)
	}
	return nil
}

func (e *executor) comparison(op *script.Comparison) error {

	err := e.addition(op.Left)
	for err == nil && op.Right != nil {
		err = e.addition(op.Right.Left)
		if err == nil {
			err = e.calculator.Op2(op.Op)
		}
		if err == nil {
			op = op.Right
		}
	}

	if err != nil {
		return errors.Error(op.Pos, err)
	}
	return nil
}

func (e *executor) addition(op *script.Addition) error {

	err := e.multiplication(op.Left)
	for err == nil && op.Right != nil {
		err = e.multiplication(op.Right.Left)
		if err == nil {
			err = e.calculator.Op2(op.Op)
		}
		if err == nil {
			op = op.Right
		}
	}

	if err != nil {
		return errors.Error(op.Pos, err)
	}
	return nil
}

func (e *executor) multiplication(op *script.Multiplication) error {

	err := e.unary(op.Left)
	for err == nil && op.Right != nil {
		err = e.unary(op.Right.Left)
		if err == nil {
			err = e.calculator.Op2(op.Op)
		}
		if err == nil {
			op = op.Right
		}
	}

	if err != nil {
		return errors.Error(op.Pos, err)
	}
	return nil
}

func (e *executor) unary(op *script.Unary) error {
	if op.Left != nil {
		err := e.primary(op.Left)
		if err == nil {
			err = e.calculator.Op1(op.Op)
		}
		if err != nil {
			return errors.Error(op.Pos, err)
		}
	}

	if op.Right != nil {
		return errors.Error(op.Pos, e.primary(op.Right))
	}

	return nil
}

func (e *executor) primary(op *script.Primary) error {
	switch {
	case op.Float != nil:
		e.calculator.Push(*op.Float)

	case op.Integer != nil:
		e.calculator.Push(*op.Integer)

	case op.String != nil:
		e.calculator.Push(*op.String)

	case op.CallFunc != nil:
		return errors.Error(op.Pos, e.callFunc(op.CallFunc))

	case op.Null, op.Nil:
		e.calculator.Push(nil)

	case op.True:
		e.calculator.Push(true)

	case op.False:
		e.calculator.Push(false)

	case op.PreIncDec != nil:
		return e.preIncDec(op.PreIncDec)

	case op.PostIncDec != nil:
		return e.postIncDec(op.PostIncDec)

	case op.Ident != nil && op.Ident.Ident != "":
		v, err := e.resolveIdent(op)
		if err != nil {
			return errors.Error(op.Pos, err)
		}

		// Just push variable onto stack
		e.calculator.Push(v)

	case op.SubExpression != nil:
		return errors.Error(op.Pos, e.Expression(op.SubExpression))

	case op.KeyValue != nil:
		if err := errors.Error(op.Pos, e.Expression(op.KeyValue.Value)); err != nil {
			return err
		}

		if val, err := e.calculator.Pop(); err != nil {
			return errors.Error(op.KeyValue.Pos, err)
		} else {
			e.calculator.Push(calculator.NewKeyValue(op.KeyValue.Key, val))
		}

	}

	return nil
}

func (e *executor) resolveIdent(op *script.Primary) (interface{}, error) {
	ident := op.Ident.Ident
	v, exists := e.state.Get(ident)
	if !exists {
		return nil, errors.Errorf(op.Pos, "%q undefined", ident)
	}

	// Handle arrays
	v, err := e.resolveArray(op, v)

	if err == nil && op.Pointer != nil {
		// Resolve references
		v, err = e.getReferenceImpl(op.Pointer, v)
	}

	if errors.IsNoFieldErr(err) {
		return v, err
	}

	if err != nil {
		return nil, errors.Error(op.Pos, err)
	}

	return v, nil
}

// preIncDec implements --i and ++i
//
// As per the C spec: --i is (i-=1) and ++i is (i+=1)
func (e *executor) preIncDec(op *script.PreIncDec) (err error) {
	// Get current value of variable
	ident := op.Ident.Ident
	value, exists := e.state.Get(ident)
	if !exists {
		return errors.Errorf(op.Pos, "%q undefined", ident)
	}

	switch {
	case op.Increment:
		value, err = calculator.Add(value, 1)
	case op.Decrement:
		value, err = calculator.Subtract(value, 1)
	default:
		// Should never occur
		err = errors.Errorf(op.Pos, "Invalid postIncDec")
	}

	if err == nil {
		// Save new value and use it
		e.state.Set(ident, value)
		e.calculator.Push(value)
	}

	return errors.Error(op.Pos, err)
}

// postIncDec implements ident++ and ident--
// where the current value of ident is returned but the variable
// is then incremented or decremented afterwards.
func (e *executor) postIncDec(op *script.PostIncDec) (err error) {
	// Get current value of variable
	ident := op.Ident.Ident
	value, exists := e.state.Get(ident)
	if !exists {
		return errors.Errorf(op.Pos, "%q undefined", ident)
	}

	newValue := value
	switch {
	case op.Increment:
		newValue, err = calculator.Add(value, 1)
	case op.Decrement:
		newValue, err = calculator.Subtract(value, 1)
	default:
		// Should never occur
		err = errors.Errorf(op.Pos, "Invalid postIncDec")
	}

	if err == nil {
		// Save new value but use the original value
		e.state.Set(ident, newValue)
		e.calculator.Push(value)
	}

	return errors.Error(op.Pos, err)
}
