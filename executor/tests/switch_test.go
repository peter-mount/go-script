package tests

import (
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/parser"
	_ "github.com/peter-mount/go-script/stdlib/fmt"
	"math"
	"testing"
)

func Test_switch(t *testing.T) {

	script1 := `main() {
  switch initial {
    case 1: result=1
    case 2: result=2
    case 3: {
      result=3
    }
	case 4.1: result=4
	case 5.1: result=5
	case 6.1: result=6
	case "A": result=65
	case "B": result=66
	case "C": result=67
	case "D", 8, "E": result=99
	case 9, "F": result=999
    default: result=42
  }
}`

	tests := []struct {
		name           string
		script         string
		params         map[string]interface{}
		initialResult  interface{}
		expectedResult interface{}
	}{
		// ===============
		// switch expression { cases... }
		// ===============
		// Integer lookup
		{
			name:           "switch exp 1",
			script:         script1,
			initialResult:  1,
			expectedResult: 1,
		},
		{
			name:           "switch exp 2",
			script:         script1,
			initialResult:  2,
			expectedResult: 2,
		},
		{
			name:           "switch exp 3",
			script:         script1,
			initialResult:  3,
			expectedResult: 3,
		},
		// float64 lookup
		{
			name:           "switch exp 4.1",
			script:         script1,
			initialResult:  4.1,
			expectedResult: 4,
		},
		{
			name:           "switch exp 5.1",
			script:         script1,
			initialResult:  5.1,
			expectedResult: 5,
		},
		{
			name:           "switch exp 6.1",
			script:         script1,
			initialResult:  6.1,
			expectedResult: 6,
		},
		// string lookup - here as this broke the parser
		{
			name:           "switch exp A",
			script:         script1,
			initialResult:  "A",
			expectedResult: 65,
		},
		{
			name:           "switch exp B",
			script:         script1,
			initialResult:  "B",
			expectedResult: 66,
		},
		{
			name:           "switch exp C",
			script:         script1,
			initialResult:  "C",
			expectedResult: 67,
		},
		{
			name:           "switch exp default",
			script:         script1,
			initialResult:  3.1415926,
			expectedResult: 42,
		},
		// multiple expressions in single case
		{
			name:           "multi-case D",
			script:         script1,
			initialResult:  "D",
			expectedResult: 99,
		},
		{
			name:           "multi-case 8",
			script:         script1,
			initialResult:  8,
			expectedResult: 99,
		},
		{
			name:           "multi-case E",
			script:         script1,
			initialResult:  "E",
			expectedResult: 99,
		},
		{
			name:           "multi-case F",
			script:         script1,
			initialResult:  "F",
			expectedResult: 999,
		},
		{
			name:           "multi-case 9",
			script:         script1,
			initialResult:  9,
			expectedResult: 999,
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

			globals.Declare("initial")
			globals.Set("initial", test.initialResult)

			globals.Declare("result")
			globals.Set("result", math.MinInt)

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
