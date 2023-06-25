package tests

import (
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/packages"
	"github.com/peter-mount/go-script/parser"
	"testing"
)

type Issue4 struct{}

// Create a Julday instance
func (_ Issue4) CreateDay() Julday {
	return 3.5
}

// Accept a Julday instance
func (_ Issue4) ProcessDay(d Julday) {}

func init() {
	packages.Register("testIssue4", &Issue4{})
}

// Test issue 4 https://github.com/peter-mount/go-script/issues/4
//
// Here we test a value returned from a function
func Test_issue_4_setFromFunc(t *testing.T) {

	tests := []struct {
		name   string
		script string
	}{
		// Test we can call a function from a global scoped variable
		{name: "setFromGlobals",
			script: `main() {
srcDay.SomeFunction()
}`},
		// Test we can call a function against a value returned by a function
		{name: "setFromFunc",
			script: `main() {
day=testIssue4.CreateDay()
day.SomeFunction()
}`},
		// Test we can pass a value to a function.
		// On the initial fix for issue 4 this broke saying not a float64
		{name: "passToFunction",
			script: `main() {
day=testIssue4.CreateDay()
testIssue4.ProcessDay(day)
}`},
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

			globals := exec.GlobalScope()
			globals.Declare("srcDay")
			globals.Set("srcDay", Julday(42))

			err = exec.Run()
			if err != nil {
				t.Fatal(err)
				return
			}
		})
	}
}
