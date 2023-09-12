package tests

import (
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/parser"
	"testing"
)

func Test_if(t *testing.T) {

	tests := []struct {
		name           string
		script         string
		params         map[string]interface{}
		initialResult  interface{}
		expectedResult interface{}
	}{
		{
			name:           "if true",
			script:         `main() { if true result=true else result=false }`,
			expectedResult: true,
		},
		{
			name:           "if false",
			script:         `main() { if false result=true else result=false }`,
			expectedResult: false,
		},
		{
			name:   "if a < b int",
			script: `main() { if a < b result=true else result=false }`,
			params: map[string]interface{}{
				"a": 1,
				"b": 2,
			},
			expectedResult: true,
		},
		{
			name:   "if a <= b int",
			script: `main() { if a <= b result=true else result=false }`,
			params: map[string]interface{}{
				"a": 1,
				"b": 2,
			},
			expectedResult: true,
		},
		{
			name:   "if a == b int",
			script: `main() { if a == b result=true else result=false }`,
			params: map[string]interface{}{
				"a": 2,
				"b": 2,
			},
			expectedResult: true,
		},
		{
			name:   "if a >= b int",
			script: `main() { if a >= b result=true else result=false }`,
			params: map[string]interface{}{
				"a": 1,
				"b": 2,
			},
			expectedResult: false,
		},
		{
			name:   "if a > b int",
			script: `main() { if a > b result=true else result=false }`,
			params: map[string]interface{}{
				"a": 1,
				"b": 2,
			},
			expectedResult: false,
		},
		{
			name:   "if a == b string true",
			script: `main() { if a == b result=true else result=false }`,
			params: map[string]interface{}{
				"a": "k4",
				"b": "k4",
			},
			expectedResult: true,
		},
		{
			name:   "if a == b string false",
			script: `main() { if a == b result=true else result=false }`,
			params: map[string]interface{}{
				"a": "k9",
				"b": "k4",
			},
			expectedResult: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			p, err := parser.New().ParseString(test.name, test.script)
			if err != nil {
				t.Fatal(err)
				return
			}

			exec, err := executor.New(p)
			if err != nil {
				t.Fatal(err)
				return
			}

			// Add each resource to the global scope
			globals := exec.GlobalScope()

			globals.Declare("result")
			globals.Set("result", test.initialResult)

			if test.params != nil {
				for k, v := range test.params {
					globals.Declare(k)
					globals.Set(k, v)
				}
			}

			// Ignore errors as we test for them
			err = exec.Run()
			if err != nil {
				t.Fatal(err)
				return
			}

			result, ok := globals.Get("result")
			if !ok {
				t.Errorf("result not returned")
			} else {
				if result != test.expectedResult {
					t.Errorf("expected %v %T got %v %T", test.expectedResult, test.expectedResult, result, result)
				}
			}

		})
	}

}
