package calculator

import "math"

var (
	operations = map[string]BiCalculation{
		"==": NewBiOpDef().
			Int(func(a, b int) (interface{}, error) { return a == b, nil }).
			Float(func(a, b float64) (interface{}, error) { return math.Abs(a-b) < 1e-9, nil }).
			String(func(a, b string) (interface{}, error) { return a == b, nil }).
			Bool(func(a, b bool) (interface{}, error) { return a == b, nil }).
			Build(),
		"!=": NewBiOpDef().
			Int(func(a, b int) (interface{}, error) { return a != b, nil }).
			Float(func(a, b float64) (interface{}, error) { return math.Abs(a-b) >= 1e-9, nil }).
			String(func(a, b string) (interface{}, error) { return a != b, nil }).
			Bool(func(a, b bool) (interface{}, error) { return a != b, nil }).
			Build(),
		"<": NewBiOpDef().
			Int(func(a, b int) (interface{}, error) { return a < b, nil }).
			Float(func(a, b float64) (interface{}, error) { return math.Abs(a-b) >= 1e-9 && a < b, nil }).
			String(func(a, b string) (interface{}, error) { return a < b, nil }).
			Build(),
		"<=": NewBiOpDef().
			Int(func(a, b int) (interface{}, error) { return a <= b, nil }).
			Float(func(a, b float64) (interface{}, error) { return math.Abs(a-b) < 1e-9 || a <= b, nil }).
			String(func(a, b string) (interface{}, error) { return a <= b, nil }).
			Build(),
		">": NewBiOpDef().
			Int(func(a, b int) (interface{}, error) { return a > b, nil }).
			Float(func(a, b float64) (interface{}, error) { return math.Abs(a-b) >= 1e-9 && a > b, nil }).
			String(func(a, b string) (interface{}, error) { return a > b, nil }).
			Build(),
		">=": NewBiOpDef().
			Int(func(a, b int) (interface{}, error) { return a >= b, nil }).
			Float(func(a, b float64) (interface{}, error) { return math.Abs(a-b) < 1e-9 || a >= b, nil }).
			String(func(a, b string) (interface{}, error) { return a >= b, nil }).
			Build(),
		"+": NewBiOpDef().
			Int(func(a, b int) (interface{}, error) { return a + b, nil }).
			Float(func(a, b float64) (interface{}, error) { return a + b, nil }).
			String(func(a, b string) (interface{}, error) { return a + b, nil }).
			Build(),
		"-": NewBiOpDef().
			Int(func(a, b int) (interface{}, error) { return a - b, nil }).
			Float(func(a, b float64) (interface{}, error) { return a - b, nil }).
			Build(),
		"*": NewBiOpDef().
			Int(func(a, b int) (interface{}, error) { return a * b, nil }).
			Float(func(a, b float64) (interface{}, error) { return a * b, nil }).
			Build(),
		"/": NewBiOpDef().
			Int(func(a, b int) (interface{}, error) { return a / b, nil }).
			Float(func(a, b float64) (interface{}, error) { return a / b, nil }).
			Build(),
	}
)
