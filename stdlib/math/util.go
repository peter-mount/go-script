package math

import (
	"github.com/peter-mount/go-script/calculator"
	"github.com/x448/float16"
)

func (_ Math) Float(v any) float64 {
	f, err := calculator.GetFloat(v)
	if err != nil {
		panic(err)
	}
	return f
}

func (m Math) Float32(v any) float32 {
	return float32(m.Float(v))
}

func (m Math) Float16(v any) float16.Float16 {
	return float16.Fromfloat32(float32(m.Float(v)))
}

func (_ Math) Int(v any) int {
	i, err := calculator.GetInt(v)
	if err != nil {
		panic(err)
	}
	return i
}
