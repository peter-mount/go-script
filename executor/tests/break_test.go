package tests

import (
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/parser"
	"testing"
)

// Test both break and continue statements
func Test_break_continue(t *testing.T) {
	tests := []struct {
		name           string
		script         string
		initialResult  int
		expectedResult int
	}{
		{
			// break after first iteration
			name:           "break 1",
			script:         `main() { for i:=1;i<10;i=i+1 { result=i break } }`,
			initialResult:  -1,
			expectedResult: 1,
		},
		{
			// break on 5th iteration after setting result
			name:           "break 5 after",
			script:         `main() { for i:=1;i<10;i=i+1 { result=i if i==5 break } }`,
			initialResult:  -1,
			expectedResult: 5,
		},
		{
			// break on 5th iteration but before setting result
			name:           "break 5 before",
			script:         `main() { for i:=1;i<10;i=i+1 { if i==5 break result=i } }`,
			initialResult:  -1,
			expectedResult: 4,
		},
		{
			// continue on first iteration but before setting result
			name:           "continue 1",
			script:         `main() { for i:=1;i<10;i=i+1 { if i>0 continue result=i } }`,
			initialResult:  -1,
			expectedResult: -1,
		},
		{
			// continue after 5th iteration
			name:           "continue 5",
			script:         `main() { for i:=0;i<10;i=i+1 { if i>5 continue result=i } }`,
			initialResult:  -1,
			expectedResult: 5,
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
					t.Errorf("expected %v got %v", test.expectedResult, result)
				}
			}

		})
	}

}
