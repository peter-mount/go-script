package tests

import (
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/parser"
	"testing"
)

func Test_assignment(t *testing.T) {
	tests := []struct {
		name           string
		script         string
		params         map[string]interface{}
		initialResult  interface{}
		expectedResult interface{}
	}{
		// ===============
		// Basic assignment
		// ===============
		{
			// Test := declare works
			name:           "basic declare",
			script:         `main() { a := 42 result=a }`,
			expectedResult: 42,
		},
		{
			// Test = works, result always exists
			name:           "basic set existing",
			script:         `main() { result= 3.1415 }`,
			expectedResult: 3.1415,
		},
		// ===============
		// assignment in inner scope
		// ===============
		{
			// declare in inner scope cannot alter outer declare
			name:           "inner declare",
			script:         `main() { a := 42 { a := 96 } result=a }`,
			expectedResult: 42,
		},
		{
			// assignment in inner scope sets outer declare
			name:           "inner set",
			script:         `main() { a := 42 { a = 96 } result=a }`,
			expectedResult: 96,
		},

		// ===============
		// Chained assignment
		//
		// a = b = c = d = f = 0
		// ===============
		{
			name:           "chained assignment declare",
			script:         `main() { a := b := c := 42 result = a==b && b==c && c==42 }`,
			expectedResult: true,
		},
		{
			name:           "chained assignment existing",
			script:         `main() { a = b = c = 42 result = a==b && b==c && c==42 }`,
			expectedResult: true,
			params:         map[string]interface{}{"a": 0, "b": 0, "c": 0},
		},
		{
			// Currently a = (b = 0 ) works so test that
			name:           "chained assignment compound",
			script:         `main() { a := (b := (c := 42)) result = a==b && b==c && c==42 }`,
			expectedResult: true,
		},

		// ===============
		// Augmented assignment
		// ===============
		{
			name:           "augmented +=",
			script:         `main() { result += 1 }`,
			initialResult:  10,
			expectedResult: 11,
		},
		{
			name:           "augmented -=",
			script:         `main() { result -= 1 }`,
			initialResult:  11,
			expectedResult: 10,
		},
		{
			name:           "augmented *=",
			script:         `main() { result *= 2 }`,
			initialResult:  10,
			expectedResult: 20,
		},
		{
			name:           "augmented /=",
			script:         `main() { result /= 2 }`,
			initialResult:  10,
			expectedResult: 5,
		},
		{
			name:           "augmented %=",
			script:         `main() { result %= 5 }`,
			initialResult:  11,
			expectedResult: 1,
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
