package calculator

import (
	"context"
	"fmt"
)

// BiCalculation performs an operation against two values
type BiCalculation interface {
	BiCalculate(a, b interface{}) (interface{}, error)
}

func CalculateOp2(ctx context.Context, op string, left, right Calculation) error {
	if err := DoCalculation(left, ctx); err == nil {
		return err
	}

	if err := DoCalculation(right, ctx); err == nil {
		return err
	}

	// If we have both then invoke the op
	if op != "" && left != nil && right != nil {
		return FromContext(ctx).Op2(op)
	}

	return nil
}

// Calculation defines an instance that can perform some calculation
type Calculation interface {
	Calculate(ctx context.Context) error
}

func DoCalculation(calc Calculation, ctx context.Context) error {
	if calc != nil {
		return calc.Calculate(ctx)
	}
	return nil
}

// BiOpDef implements an operation whose behaviour depends on the type
// of the left hand side
type BiOpDef struct {
	intOp    func(a, b int) (interface{}, error)
	floatOp  func(a, b float64) (interface{}, error)
	stringOp func(a, b string) (interface{}, error)
	boolOp   func(a, b bool) (interface{}, error)
}

func (op *BiOpDef) doInt(a int, b interface{}) (interface{}, error) {
	c, err := GetInt(b)
	if err != nil {
		return nil, err
	}
	return op.intOp(a, c)
}

func (op *BiOpDef) doFloat(a float64, b interface{}) (interface{}, error) {
	c, err := GetFloat(b)
	if err != nil {
		return nil, err
	}
	return op.floatOp(a, c)
}

func (op *BiOpDef) doString(a string, b interface{}) (interface{}, error) {
	c, err := GetString(b)
	if err != nil {
		return nil, err
	}
	return op.stringOp(a, c)
}

func (op *BiOpDef) doBool(a bool, b interface{}) (interface{}, error) {
	c, err := GetBool(b)
	if err != nil {
		return nil, err
	}
	return op.boolOp(a, c)
}

func (op *BiOpDef) BiCalculate(a, b interface{}) (interface{}, error) {
	if op.intOp != nil {
		if i, ok := a.(int); ok {
			return op.doInt(i, b)
		}
		if i, ok := a.(int64); ok {
			return op.doInt(int(i), b)
		}
		if i, ok := a.(Int); ok {
			return op.doInt(i.Int(), b)
		}
	}

	if op.floatOp != nil {
		if f, ok := a.(float64); ok {
			return op.doFloat(f, b)
		}
		if f, ok := a.(Float); ok {
			return op.doFloat(f.Float(), b)
		}
		if f, ok := a.(float32); ok {
			return op.doFloat(float64(f), b)
		}
	}

	if op.stringOp != nil {
		if s, ok := a.(string); ok {
			return op.doString(s, b)
		}
		if s, ok := a.(String); ok {
			return op.doString(s.String(), b)
		}
	}

	if op.boolOp != nil {
		if ab, ok := a.(bool); ok {
			return op.doBool(ab, b)
		}
	}

	return nil, fmt.Errorf("unable to convert %T to %T", b, a)
}

// NewBiOpDef creates a new NewBiOp based on the provided
func NewBiOpDef() BiOpDefBuilder {
	return &biOpBuilder{}
}

type BiOpDefBuilder interface {
	Int(func(a, b int) (interface{}, error)) BiOpDefBuilder
	Float(func(a, b float64) (interface{}, error)) BiOpDefBuilder
	String(func(a, b string) (interface{}, error)) BiOpDefBuilder
	Bool(func(a, b bool) (interface{}, error)) BiOpDefBuilder
	Build() *BiOpDef
}

type biOpBuilder struct {
	fInt    func(a, b int) (interface{}, error)
	fFloat  func(a, b float64) (interface{}, error)
	fString func(a, b string) (interface{}, error)
	fBool   func(a, b bool) (interface{}, error)
}

func (b *biOpBuilder) Int(f func(a, b int) (interface{}, error)) BiOpDefBuilder {
	b.fInt = f
	return b
}

func (b *biOpBuilder) Float(f func(a, b float64) (interface{}, error)) BiOpDefBuilder {
	b.fFloat = f
	return b
}

func (b *biOpBuilder) String(f func(a, b string) (interface{}, error)) BiOpDefBuilder {
	b.fString = f
	return b
}

func (b *biOpBuilder) Bool(f func(a, b bool) (interface{}, error)) BiOpDefBuilder {
	b.fBool = f
	return b
}

func (b *biOpBuilder) Build() *BiOpDef {
	return &BiOpDef{
		intOp:    b.fInt,
		floatOp:  b.fFloat,
		stringOp: b.fString,
		boolOp:   b.fBool,
	}
}
