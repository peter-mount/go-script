package executor

import (
	"github.com/peter-mount/go-script/script"
	"reflect"
	"strconv"
	"testing"
)

// Test_executor_callReflectFuncImpl tests the inner workings of reflective function calls.
// Main issue here is testing variadic function calls
func Test_executor_callReflectFuncImpl(t *testing.T) {
	// Dummies needed to make calls
	e := &executor{}
	cf := &script.CallFunc{}

	tests := []struct {
		// Arguments to pass to function
		args []interface{}
		// Wanted return value(s)
		want interface{}
		// true if we want it to fail
		wantErr bool
	}{
		// Differing number of var args including none
		{[]interface{}{1.0, 2.0, 3.0}, 6.0, false},
		{[]interface{}{1.0, 2.0, 3.0, 4.0}, 10.0, false},
		{[]interface{}{1.0, 2.0, 3.0, 4.0, 5.0}, 15.0, false},
		{[]interface{}{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}, 21.0, false},
		{[]interface{}{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0}, 28.0, false},
		// Here we are passing an empty variadic so only the fixed parameters are present
		{[]interface{}{1.0, 2.0}, 3.0, false},
		// Special case, passing different values in the variadic section
		{[]interface{}{1, 2, 3, 4}, 10.0, false},
		{[]interface{}{1.0, 2.0, 3.0, 4, 5.0, "6", 7.0}, 28.0, false},
	}

	for ti, tt := range tests {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			// Reset run test
			run := false
			variadic := func(a, b float64, c ...float64) float64 {
				// Mark we have been called
				run = true

				// Test we have the correct number of arguments passed to the function
				// -2 is to account for a & b
				if len(c) != (len(tt.args) - 2) {
					t.Errorf("variadic expected %d elements got %d ", len(tt.args)-2, len(c))
				}

				// Calculate sum of all arguments
				total := a + b
				for _, d := range c {
					total = total + d
				}
				return total
			}

			got, err := e.callReflectFuncImpl(cf, reflect.ValueOf(variadic), tt.args)
			if err != nil {
				t.Error(err)
				return
			}

			if tt.wantErr {
				t.Errorf("Expected error to be returned")
			}

			if !run {
				t.Errorf("Variadic did not run")
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Got %T %v wanted %T %v", got, got, tt.want, tt.want)
			}
		})
	}
}
