package tests

import (
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/parser"
	"testing"
)

func Test_ternary(t *testing.T) {

	tests := []struct {
		name           string
		script         string
		params         map[string]interface{}
		initialResult  interface{}
		expectedResult interface{}
	}{
		{
			name:           "if true",
			script:         `main() { result = true ? true : false }`,
			expectedResult: true,
		},
		{
			name:           "if false",
			script:         `main() { result = false ? true : false }`,
			expectedResult: false,
		},
		{
			name:   "if a<b",
			script: `main() { result = a<b ? 10 : 30 }`,
			params: map[string]interface{}{
				"a": 1,
				"b": 5,
			},
			expectedResult: 10,
		},
		{
			name:   "if a>b",
			script: `main() { result = a>b ? 10 : 20 }`,
			params: map[string]interface{}{
				"a": 1,
				"b": 5,
			},
			expectedResult: 20,
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
