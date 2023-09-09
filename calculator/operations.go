package calculator

import "math"

var (
	monoOperations = map[string]MonoCalculation{
		"!": NewMonoOpDef().
			Int(func(a int) (interface{}, error) { return a == 0, nil }).
			Float(func(a float64) (interface{}, error) { return math.Abs(a) <= 1e-9, nil }).
			Bool(func(a bool) (interface{}, error) { return !a, nil }).
			Build(),
		"-": NewMonoOpDef().
			Int(func(a int) (interface{}, error) { return -a, nil }).
			Float(func(a float64) (interface{}, error) { return -a, nil }).
			Bool(func(a bool) (interface{}, error) { return !a, nil }).
			Build(),
	}

	biOperations = map[string]BiCalculation{
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
		"&&": NewBiOpDef().
			Bool(func(a, b bool) (interface{}, error) { return a && b, nil }).
			Build(),
		"||": NewBiOpDef().
			Bool(func(a, b bool) (interface{}, error) { return a || b, nil }).
			Build(),
		"%": NewBiOpDef().
			Int(func(a, b int) (interface{}, error) { return a % b, nil }).
			Float(func(a, b float64) (interface{}, error) { return math.Mod(a, b), nil }).
			Build(),
	}
)
