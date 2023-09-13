package calculator

import (
	"fmt"
)

// MonoCalculation performs an operation against two values
type MonoCalculation interface {
	MonoCalculate(a interface{}) (interface{}, error)
}

// BiCalculation performs an operation against two values
type BiCalculation interface {
	BiCalculate(a, b interface{}) (interface{}, error)
}

type biOpCommon struct {
	intOp    func(a, b int) (interface{}, error)
	floatOp  func(a, b float64) (interface{}, error)
	stringOp func(a, b string) (interface{}, error)
	boolOp   func(a, b bool) (interface{}, error)
}

// BiOpDef implements an operation whose behaviour depends on the type
// of the left hand side
type BiOpDef struct {
	biOpCommon
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

	// Convert so a and b are the same type
	a, b, err := Convert(a, b)
	if err != nil {
		return nil, err
	}

	if op.floatOp != nil {
		if f, ok := GetFloatRaw(a); ok {
			return op.doFloat(f, b)
		}
	}

	if op.intOp != nil {
		if i, ok := GetIntRaw(a); ok {
			return op.doInt(i, b)
		}
	}

	if op.stringOp != nil {
		if s, ok := GetStringRaw(a); ok {
			return op.doString(s, b)
		}
	}

	if op.boolOp != nil {
		if ab, ok := GetBoolRaw(a); ok {
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
	biOpCommon
}

func (b *biOpBuilder) Int(f func(a, b int) (interface{}, error)) BiOpDefBuilder {
	b.intOp = f
	return b
}

func (b *biOpBuilder) Float(f func(a, b float64) (interface{}, error)) BiOpDefBuilder {
	b.floatOp = f
	return b
}

func (b *biOpBuilder) String(f func(a, b string) (interface{}, error)) BiOpDefBuilder {
	b.stringOp = f
	return b
}

func (b *biOpBuilder) Bool(f func(a, b bool) (interface{}, error)) BiOpDefBuilder {
	b.boolOp = f
	return b
}

func (b *biOpBuilder) Build() *BiOpDef {
	return &BiOpDef{biOpCommon: b.biOpCommon}
}

// NewMonoOpDef creates a new NewBiOp based on the provided
func NewMonoOpDef() MonoOpDefBuilder {
	return &monoOpBuilder{}
}

type MonoOpDefBuilder interface {
	Int(func(a int) (interface{}, error)) MonoOpDefBuilder
	Float(func(a float64) (interface{}, error)) MonoOpDefBuilder
	String(func(a string) (interface{}, error)) MonoOpDefBuilder
	Bool(func(a bool) (interface{}, error)) MonoOpDefBuilder
	Build() *MonoOpDef
}

type monoOpCommon struct {
	intOp    func(a int) (interface{}, error)
	floatOp  func(a float64) (interface{}, error)
	stringOp func(a string) (interface{}, error)
	boolOp   func(a bool) (interface{}, error)
}

// MonoOpDef implements an operation whose behaviour depends on the type
// of the left hand side
type MonoOpDef struct {
	monoOpCommon
}

func (op *MonoOpDef) MonoCalculate(a interface{}) (interface{}, error) {

	if op.floatOp != nil {
		if f, ok := GetFloatRaw(a); ok {
			return op.floatOp(f)
		}
	}

	if op.intOp != nil {
		if i, ok := GetIntRaw(a); ok {
			return op.intOp(i)
		}
	}

	if op.stringOp != nil {
		if s, ok := GetStringRaw(a); ok {
			return op.stringOp(s)
		}
	}

	if op.boolOp != nil {
		if ab, ok := GetBoolRaw(a); ok {
			return op.boolOp(ab)
		}
	}

	return nil, fmt.Errorf("unable to convert %T", a)
}

type monoOpBuilder struct {
	monoOpCommon
}

func (b *monoOpBuilder) Int(f func(a int) (interface{}, error)) MonoOpDefBuilder {
	b.intOp = f
	return b
}

func (b *monoOpBuilder) Float(f func(a float64) (interface{}, error)) MonoOpDefBuilder {
	b.floatOp = f
	return b
}

func (b *monoOpBuilder) String(f func(a string) (interface{}, error)) MonoOpDefBuilder {
	b.stringOp = f
	return b
}

func (b *monoOpBuilder) Bool(f func(a bool) (interface{}, error)) MonoOpDefBuilder {
	b.boolOp = f
	return b
}

func (b *monoOpBuilder) Build() *MonoOpDef {
	return &MonoOpDef{monoOpCommon: b.monoOpCommon}
}
