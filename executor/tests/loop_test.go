package tests

import (
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/parser"
	"testing"
)

// Test both break and continue statements
func Test_loops(t *testing.T) {
	tests := []struct {
		name           string
		script         string
		initialResult  int
		expectedResult int
	}{
		// ===============
		// FOR
		// ===============
		{
			// Simple loop, result is value that failed the condition
			name:           "for simple",
			script:         `main() { for i:=1;i<10;i=i+1 { result=i } }`,
			initialResult:  -1,
			expectedResult: 9,
		},
		{
			// Increment result with empty body
			name:           "for no-init {}",
			script:         `main() { for ;result<10;result=result+1 { } }`,
			initialResult:  0,
			expectedResult: 10,
		},
		{
			// Increment result in body
			name:           "for no-init-inc {}",
			script:         `main() { for ;result<10; { result=result+1 } }`,
			initialResult:  0,
			expectedResult: 10,
		},
		{
			// Increment result in body test in body
			name:           "for no-init-inc break{}",
			script:         `main() { for ;; { result=result+1 if result > 9 break } }`,
			initialResult:  0,
			expectedResult: 10,
		},
		{
			// Increment result with empty ; statement
			name:           "for no-init ;",
			script:         `main() { for ;result<10;result=result+1; }`,
			initialResult:  0,
			expectedResult: 10,
		},
		{
			// Should never run
			name:           "for no-init never run",
			script:         `main() { for ;result<0;result=result+1 { } }`,
			initialResult:  0,
			expectedResult: 0,
		},

		// ===============
		// FOR with post increment/decrement
		// ===============
		{
			name:           "for i++",
			script:         `main() { for i:=1; i<10; i++ { result=i } }`,
			initialResult:  -1,
			expectedResult: 9,
		},
		{
			name:           "for i--",
			script:         `main() { for i:=10; i>0; i-- { result=i } }`,
			initialResult:  10,
			expectedResult: 1,
		},

		// ===============
		// FOR with pre increment/decrement
		// ===============
		{
			name:           "for ++i",
			script:         `main() { for i:=1; i<10; ++i { result=i } }`,
			initialResult:  -1,
			expectedResult: 9,
		},
		{
			name:           "for --i",
			script:         `main() { for i:=10; i>0; --i { result=i } }`,
			initialResult:  10,
			expectedResult: 1,
		},

		// ===============
		// DO WHILE
		// ===============
		{
			// Simple loop, result is value that failed the condition
			name:           "do while",
			script:         `main() { do result=result+1 while result < 20 }`,
			initialResult:  0,
			expectedResult: 20,
		},
		{
			// Test we always run at least once
			name:           "do while always run",
			script:         `main() { do result=result+1 while result < 0 }`,
			initialResult:  10,
			expectedResult: 11,
		},

		// ===============
		// REPEAT UNTIL
		// ===============
		{
			// Simple loop, result is value that failed the condition
			name:           "repeat until",
			script:         `main() { repeat result=result+1 until result >= 20 }`,
			initialResult:  0,
			expectedResult: 20,
		},
		{
			// Test we always run at least once
			name:           "repeat until always run",
			script:         `main() { repeat result=result+1 until result > 0 }`,
			initialResult:  10,
			expectedResult: 11,
		},

		// ===============
		// WHILE
		// ===============
		{
			// Simple loop, result is value that failed the condition
			name:           "while",
			script:         `main() { while result < 15 result=result+1 }`,
			initialResult:  0,
			expectedResult: 15,
		},
		{
			// Test we never run if the condition was always false
			name:           "while never run",
			script:         `main() { while result < 15 result=result+1 }`,
			initialResult:  15,
			expectedResult: 15,
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
