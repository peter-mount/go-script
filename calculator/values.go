package calculator

import (
	"fmt"
	"math"
	"strconv"
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

// Bool is an instance that can return its value as a bool
type Bool interface {
	Bool() bool
}

// GetFloatRaw returns v as a float64 if it's a form of float.
func GetFloatRaw(v interface{}) (float64, bool) {
	if v != nil {
		if f, ok := v.(float64); ok {
			return f, true
		}
		if f, ok := v.(Float); ok {
			return f.Float(), true
		}
		if f, ok := v.(float32); ok {
			return float64(f), true
		}
	}
	return 0, false
}

// GetFloat returns v as a float64.
// If v is a float64 it will return it.
// If v is an int or int64 it will return that value as a float64.
// If v implements the Float interface then it will use that for the result.
// If v is a string it will parse it.
// The bool is false if a float64 cannot be returned.
func GetFloat(v interface{}) (float64, error) {
	if f, ok := GetFloatRaw(v); ok {
		return f, nil
	}

	if v != nil {
		if i, ok := GetIntRaw(v); ok {
			return float64(i), nil
		}

		if s, ok := v.(string); ok {
			return strconv.ParseFloat(s, 64)
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

// GetIntRaw returns v as an int if v is a form of integer.
func GetIntRaw(v interface{}) (int, bool) {
	if v != nil {
		if i, ok := v.(int); ok {
			return i, true
		}
		if i, ok := v.(Int); ok {
			return i.Int(), true
		}
		if i, ok := v.(int64); ok {
			return int(i), true
		}
	}
	return 0, false
}

// GetInt returns v as an int.
// If v is an int or int64 it will return that value.
// If v implements the Int interface then it will use that for the result.
// If v is a float64 it will return it as an int.
// If v is a string it will parse it.
// The bool is false if an integer cannot be returned.
func GetInt(v interface{}) (int, error) {
	if i, ok := GetIntRaw(v); ok {
		return i, nil
	}

	if v != nil {
		if f, ok := GetFloatRaw(v); ok {
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
	}
	return 0, fmt.Errorf("not an int %q", v)
}

// GetStringRaw returns v as a string if it's a string or implements String.
// Returns "",false if the value is not a string.
func GetStringRaw(v interface{}) (string, bool) {
	if v != nil {
		if s, ok := v.(string); ok {
			return s, true
		}

		if s, ok := v.(String); ok {
			return s.String(), true
		}
	}

	return "", false
}

// GetString returns v as a string.
// If v is string or implements the String interface then it will use
// that result.
// If an int, int64, float64 then it will return that value as a string.
// If a bool then "true" or "false" is returned.
// Returns "",false if the value could not be converted to a string.
func GetString(v interface{}) (string, error) {
	if s, ok := GetStringRaw(v); ok {
		return s, nil
	}

	if v != nil {
		if i, ok := GetIntRaw(v); ok {
			return strconv.Itoa(i), nil
		}

		if f, ok := GetFloatRaw(v); ok {
			return strconv.FormatFloat(f, 'f', 6, 64), nil
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

// GetBoolRaw returns a bool if v is a bool or implements Bool
func GetBoolRaw(v interface{}) (bool, bool) {
	if v != nil {
		if b, ok := v.(bool); ok {
			return b, true
		}
		if b, ok := v.(Bool); ok {
			return b.Bool(), true
		}
	}
	return false, false
}

// GetBool converts v to a bool, returning an error if it cannot do the conversion.
//
// For int this returns true if it's not 0.
//
// For float this returns true if |float| > 1e-9 (to account for rounding errors)
//
// For string this returns true if "true", "yes", "t" or "y", and false if "false", "no", "f" or "n".
func GetBool(v interface{}) (bool, error) {
	if b, ok := GetBoolRaw(v); ok {
		return b, nil
	}

	if v != nil {
		if s, ok := GetStringRaw(v); ok {
			switch s {
			case "true", "yes", "t", "y":
				return true, nil
			case "false", "no", "f", "n":
				return false, nil
			}
		}

		if i, ok := GetIntRaw(v); ok {
			return i != 0, nil
		}

		if f, ok := GetFloatRaw(v); ok {
			return math.Abs(f) >= 1e-9, nil
		}

	}

	return false, fmt.Errorf("not a bool %q", v)
}

// Convert converts 'a' and 'b' so that they are of the same type.
// e.g. if 'a' is a float then this will ensure both are float's.
// Same for int, string or bool.
//
// If the conversion cannot take place then this will return an error.
//
// Special case: If 'a' is an int but 'b' is a float then this will
// convert 'a' to a float.
//
// This is to allow "for i=0; i<10; i=i+0.5" to work because in the
// increment "i=i+0.5" "i" is an int and if we convert 0.5 to an int
// then we get 0 and an infinite loop.
func Convert(a, b interface{}) (interface{}, interface{}, error) {
	af, aFloat := GetFloatRaw(a)
	bf, bFloat := GetFloatRaw(b)
	ai, aInt := GetIntRaw(a)
	switch {

	// a & b are floats so leave alone
	case aFloat && bFloat:
		return af, bf, nil

	// 'a' is int, 'b' is float so return floats
	case aInt && bFloat:
		return float64(ai), bf, nil

	// 'a' float so try to convert 'b' to float
	case aFloat:
		f, err := GetFloat(b)
		return af, f, err

	// 'a' is int so try to convert 'b' to int
	case aInt:
		i, err := GetInt(b)
		return ai, i, err

	default:
		if as, ok := GetStringRaw(a); ok {
			bs, err := GetString(b)
			return as, bs, err
		}

		if ab, ok := GetBoolRaw(a); ok {
			bb, err := GetBool(b)
			return ab, bb, err
		}
	}

	return nil, nil, fmt.Errorf("unable to convert %T to %T", b, a)
}
