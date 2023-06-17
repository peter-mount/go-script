package math

import "github.com/peter-mount/go-script/calculator"

func (_ Math) Float(v any) float64 {
	f, err := calculator.GetFloat(v)
	if err != nil {
		panic(err)
	}
	return f
}

func (_ Math) Int(v any) int {
	i, err := calculator.GetInt(v)
	if err != nil {
		panic(err)
	}
	return i
}
