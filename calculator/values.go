package calculator

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Float is an instance that can return its value as a float64
type Float interface {
	Float() float64
}

// Int is an instance that can return its value as an int
type Int interface {
	Int() int
}

// String is an instance that can return its value as a string
type String interface {
	String() string
}

// GetFloat returns v as a float64.
// If v is a float64 it will return it.
// If v is an int or int64 it will return that value as a float64.
// If v implements the Float interface then it will use that for the result.
// If v is a string it will parse it.
// The bool is false if a float64 cannot be returned.
func GetFloat(v interface{}) (float64, error) {
	if v != nil {
		if f, ok := v.(float64); ok {
			return f, nil
		}
		if f, ok := v.(Float); ok {
			return f.Float(), nil
		}
		if f, ok := v.(float32); ok {
			return float64(f), nil
		}
		if s, ok := v.(string); ok {
			return strconv.ParseFloat(s, 64)
		}
		if f, ok := v.(int); ok {
			return float64(f), nil
		}
		if f, ok := v.(Int); ok {
			return float64(f.Int()), nil
		}
		if f, ok := v.(int64); ok {
			return float64(f), nil
		}
		if b, ok := v.(bool); ok {
			if b {
				return 1, nil
			}
			return 0, nil
		}
	}
	return 0, fmt.Errorf("not float %q", v)
}

// GetInt returns v as an int.
// If v is an int or int64 it will return that value.
// If v implements the Int interface then it will use that for the result.
// If v is a float64 it will return it as an int.
// If v is a string it will parse it.
// The bool is false if an integer cannot be returned.
func GetInt(v interface{}) (int, error) {
	if v != nil {
		if f, ok := v.(int); ok {
			return f, nil
		}
		if f, ok := v.(Int); ok {
			return f.Int(), nil
		}
		if f, ok := v.(int64); ok {
			return int(f), nil
		}
		if s, ok := v.(string); ok {
			return strconv.Atoi(s)
		}
		if b, ok := v.(bool); ok {
			if b {
				return 1, nil
			}
			return 0, nil
		}
		if f, ok := v.(float64); ok {
			return int(f), nil
		}
		if f, ok := v.(Float); ok {
			return int(f.Float()), nil
		}
		if f, ok := v.(float32); ok {
			return int(f), nil
		}
	}
	return 0, fmt.Errorf("not an int %q", v)
}

// GetString returns v as a string.
// If v is string or implements the String interface then it will use
// that result.
// If an int, int64, float64 then it will return that value as a string.
// Returns "",false if the value could not be converted to a string.
func GetString(v interface{}) (string, error) {
	if v != nil {
		if s, ok := v.(string); ok {
			return s, nil
		}

		if s, ok := v.(String); ok {
			return s.String(), nil
		}

		if f, ok := v.(int); ok {
			return strconv.Itoa(f), nil
		}
		if f, ok := v.(Int); ok {
			return strconv.Itoa(f.Int()), nil
		}
		if f, ok := v.(int64); ok {
			return strconv.Itoa(int(f)), nil
		}

		if f, ok := v.(float64); ok {
			return strconv.FormatFloat(f, 'f', 6, 64), nil
		}
		if f, ok := v.(Float); ok {
			return strconv.FormatFloat(f.Float(), 'f', 6, 64), nil
		}
		if f, ok := v.(float32); ok {
			return strconv.FormatFloat(float64(f), 'f', 6, 32), nil
		}

		if b, ok := v.(bool); ok {
			if b {
				return "true", nil
			}
			return "false", nil
		}
	}

	return "", fmt.Errorf("not a string %q", v)
}

func GetBool(v interface{}) (bool, error) {
	if v != nil {
		if b, ok := v.(bool); ok {
			return b, nil
		}

		if s, ok := v.(String); ok {
			v = s.String()
		}
		if s, ok := v.(string); ok {
			switch s {
			case "true", "yes", "t", "y":
				return true, nil
			case "false", "no", "f", "n":
				return false, nil
			}
		}

		if f, ok := v.(int); ok {
			return f != 0, nil
		}
		if f, ok := v.(Int); ok {
			return f.Int() != 0, nil
		}
		if f, ok := v.(int64); ok {
			return f != 0, nil
		}

		if f, ok := v.(float64); ok {
			return math.Abs(f) < 1e-9, nil
		}
		if f, ok := v.(Float); ok {
			return math.Abs(f.Float()) < 1e-9, nil
		}
		if f, ok := v.(float32); ok {
			return math.Abs(float64(f)) < 1e-9, nil
		}

	}

	return false, fmt.Errorf("not a bool %q", v)
}

// Convert converts b so that it's of the same type as a
func Convert(a, b interface{}) (interface{}, error) {
	if _, ok := a.(float64); ok {
		return GetFloat(b)
	}
	if _, ok := a.(Float); ok {
		return GetFloat(b)
	}
	if _, ok := a.(float32); ok {
		return GetFloat(b)
	}

	if _, ok := a.(int); ok {
		return GetInt(b)
	}
	if _, ok := a.(Int); ok {
		return GetInt(b)
	}
	if _, ok := a.(int64); ok {
		return GetInt(b)
	}

	if _, ok := a.(string); ok {
		return GetString(b)
	}
	if _, ok := a.(String); ok {
		return GetString(b)
	}

	return nil, fmt.Errorf("unable to convert %T to %T", b, a)
}

func GetValue(t string, v interface{}) (interface{}, error) {
	switch strings.ToLower(t) {
	case "", "void":
		return nil, nil
	case "int":
		return GetInt(v)

	case "float":
		return GetFloat(v)

	case "string":
		return GetString(v)

	case "bool":
		return GetBool(v)

	default:
		return nil, fmt.Errorf("type %q is unsupported", t)
	}
}
