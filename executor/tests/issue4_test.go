package tests

import (
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/packages"
	"github.com/peter-mount/go-script/parser"
	"testing"
)

type Issue4 struct{}

func (_ Issue4) CreateDay() Julday {
	return 3.5
}

func init() {
	packages.Register("testIssue4", &Issue4{})
}

// Test issue 4 https://github.com/peter-mount/go-script/issues/4
//
// Here we test that the value does work if set directly in a variable
func Test_issue_4_setFromGlobals(t *testing.T) {

	script := `main() {
day.SomeFunction()
}`
	p, err := parser.New().ParseString("", script)
	if err != nil {
		t.Fatal(err)
		return
	}

	exec, err := executor.New(p)
	if err != nil {
		t.Fatal(err)
		return
	}

	gs := exec.GlobalScope()
	gs.Declare("day")
	gs.Set("day", Julday(3.5))

	err = exec.Run()
	if err != nil {
		t.Fatal(err)
		return
	}
}

// Test issue 4 https://github.com/peter-mount/go-script/issues/4
//
// Here we test a value returned from a function
func Test_issue_4_setFromFunc(t *testing.T) {

	script := `main() {
day=testIssue4.CreateDay()
day.SomeFunction()
}`
	p, err := parser.New().ParseString("", script)
	if err != nil {
		t.Fatal(err)
		return
	}

	exec, err := executor.New(p)
	if err != nil {
		t.Fatal(err)
		return
	}

	err = exec.Run()
	if err != nil {
		t.Fatal(err)
		return
	}
}
