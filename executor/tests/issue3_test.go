package tests

import (
	"reflect"
	"testing"
)

type Issue3 struct {
	Value int
}

// NonPointerFunction should always work
func (s Issue3) NonPointerFunction() bool { return true }

// PointerFunction should work but does not if Issue3 is not passed as a pointer
func (s *Issue3) PointerFunction() bool { return true }

// Test issue 3 https://github.com/peter-mount/go-script/issues/3
func Test_issue_3(t *testing.T) {

	tests := []struct {
		Pointer  bool
		Value    any
		Function string
		Valid    bool
	}{
		{Pointer: true, Value: &Issue3{Value: 42}, Function: "NonPointerFunction", Valid: true},
		{Pointer: true, Value: &Issue3{Value: 42}, Function: "PointerFunction", Valid: true},
		{Value: Issue3{Value: 42}, Function: "NonPointerFunction", Valid: true},
		// This should return a valid method lookup, but it does not.
		// It seems that this is not possible with go's reflection.
		// If the struct is passed to reflect.ValueOf() is not a pointer then
		// it's pointer methods are not visible and there is no possible way
		// to get a pointer to it.
		{Value: Issue3{Value: 42}, Function: "PointerFunction", Valid: false},
	}

	for _, test := range tests {
		n := test.Function
		if test.Pointer {
			n = "Pointer " + n
		} else {
			n = "NonPointer " + n
		}
		t.Run(n, func(t *testing.T) {
			ti := reflect.ValueOf(test.Value)

			tf := ti.MethodByName(test.Function)

			if test.Valid != tf.IsValid() {
				t.Errorf("Failed lookup %q expected valid=%v got %v", test.Function, test.Valid, tf.IsValid())
			}

			if tf.IsValid() {
				retV := tf.Call([]reflect.Value{})
				if len(retV) == 0 {
					t.Errorf("Return slice empty")
				}
				for _, r := range retV {
					if !r.IsValid() {
						t.Errorf("Return value not valid")
					}

					if r.Kind() != reflect.Bool {
						t.Errorf("Return value not Bool")
					}

					if !r.Bool() {
						t.Errorf("Return not true")
					}
					return
				}
			}

		})

	}
}
