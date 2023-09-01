package executor

import (
	"fmt"
	"github.com/peter-mount/go-script/script"
	"reflect"
	"testing"
)

// Test_executor_callReflectFuncImpl tests the inner workings of reflective function calls.
// Main issue here is testing variadic function calls
func Test_executor_callReflectFuncImpl(t *testing.T) {
	// Dummies needed to make calls
	e := &executor{}
	cf := &script.CallFunc{
		// 3 expressions for 3 args in test function
		Parameters: &script.ParameterList{
			Args: []*script.Expression{{}, {}, {}},
		},
	}

	tests := []struct {
		// Arguments to pass to function
		args []interface{}
		// Wanted return value(s)
		want interface{}
		// true if we want it to fail
		wantErr bool
		// True to mark the function call as variadic
		variadic bool
	}{
		// Differing number of var args including none
		{args: []interface{}{1.0, 2.0, 3.0}, want: 6.0},
		{args: []interface{}{1.0, 2.0, 3.0, 4.0}, want: 10.0},
		{args: []interface{}{1.0, 2.0, 3.0, 4.0, 5.0}, want: 15.0},
		{args: []interface{}{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}, want: 21.0},
		{args: []interface{}{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0}, want: 28.0},
		// Here we are passing an empty variadic so only the fixed parameters are present
		{args: []interface{}{1.0, 2.0}, want: 3.0},
		// Special case, passing different values in the variadic section
		{args: []interface{}{1, 2, 3, 4}, want: 10.0},
		{args: []interface{}{1.0, 2.0, 3.0, 4, 5.0, "6", 7.0}, want: 28.0},
		// Test variadic function calls, here it should expand the slice
		{args: []interface{}{1.0, 2.0, []interface{}{3.0, 4.0, 5.0, 6.0, 7.0}}, want: 28.0, variadic: true},
		{args: []interface{}{1.0, 2.0, 3.0, []interface{}{4.0, 5.0, 6.0, 7.0}}, want: 28.0, variadic: true},
		{args: []interface{}{1.0, 2.0, 3.0, 4.0, []interface{}{5.0, 6.0, 7.0}}, want: 28.0, variadic: true},
		{args: []interface{}{1.0, 2.0, 3.0, 4.0, 5.0, []interface{}{6.0, 7.0}}, want: 28.0, variadic: true},
		// Variadic call but last entry is not a slice - this should pass as here this is valid
		{args: []interface{}{1.0, 2.0, 3.0, 4, 5.0, "6", 7.0}, want: 28.0, variadic: true},
		// Should fail as we pass a slice, but it is variadic
		{args: []interface{}{1.0, 2.0, []interface{}{3.0, 4.0, 5.0, 6.0, 7.0}}, want: 28.0, wantErr: true},
		// Should fail as last entry is not a slice but a preceding one is
		{args: []interface{}{1.0, 2.0, []interface{}{3.0, 4.0, 5.0, 6.0}, 7.0}, want: 28.0, wantErr: true},
		{args: []interface{}{1.0, 2.0, []interface{}{3.0, 4.0, 5.0}, 6.0, 7.0}, want: 28.0, wantErr: true},
	}

	for ti, tt := range tests {
		t.Run(fmt.Sprintf("%02d %v %v %v", ti, tt.want, tt.variadic, tt.wantErr),
			func(t *testing.T) {
				// Mark the call as Variadic - e.g. '...' after last argument
				if cf.Parameters != nil {
					cf.Parameters.Variadic = tt.variadic
				}

				// Reset run test
				run := false
				variadic := func(a, b float64, c ...float64) float64 {
					// Mark we have been called
					run = true

					// Calculate sum of all arguments
					total := a + b
					for _, d := range c {
						total = total + d
					}
					return total
				}

				got, err := e.CallReflectFuncImpl(cf, reflect.ValueOf(variadic), tt.args)
				if err != nil {
					if !tt.wantErr {
						t.Error(err)
					}
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
