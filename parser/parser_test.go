package parser

import (
	"github.com/peter-mount/go-script/executor"
	"strings"
	"testing"
)

// Test_parser deals with parser issues
//
// e.g. animUtil.Rect(dw[0]*(dw[1]-1), lb.Y1, lb.X2, lb.Y2) failed to parse
// after refactoring as it expected ++ or -- when it encountered * in the above line.
func Test_parser(t *testing.T) {
	tests := []struct {
		name          string
		script        string
		expectedError string
	}{
		// ====================================================================
		// After a refactor, dw[0]*(dw[1]-1) failed to parse with:
		//
		// unexpected token "*" (expected (("-" "-") | ("+" "+")))
		//
		// If the [0] and [1] are removed so just referencing dw then it parses
		// ====================================================================
		{
			name:   "expression dw[0]*dw[1]",
			script: `main() { dw[0]*(dw[1]-1) }`,
			//script: `main() { animUtil.Rect(dw[0]*(dw[1]-1), lb.Y1, lb.X2, lb.Y2) }`,
		},
		{
			name:   "expression dw*dw",
			script: `main() { dw*(dw-1) }`,
		},

		// ====================================================================
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			p, err := New().ParseString(test.name, test.script)
			if test.expectedError != "" {
				if err == nil {
					t.Fatalf("expected %q but got no error", test.expectedError)
				} else {
					msg := err.Error()
					if !strings.Contains(msg, test.expectedError) {
						t.Fatalf("expected %q but got %q", test.expectedError, msg)
					}
				}
				// We stop the test here as we got the expectedError or not
				return
			} else if err != nil {
				t.Fatal(err)
				return
			}

			// Create the executor, so we also test the parsed script
			// initialises correctly as well
			_, err = executor.New(p)
			if err != nil {
				t.Fatal(err)
				return
			}

		})
	}

}
