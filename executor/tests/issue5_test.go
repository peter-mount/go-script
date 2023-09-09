package tests

import (
	"fmt"
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/parser"
	"testing"
)

type issue5resource struct {
	Id           int
	Counter      *int // source of Sequence
	CreateCloser bool // true if this is a CreateCloser
	CloseError   bool // true to return an error on close
	CreateError  bool // true to return an error on create
	issue5results
}

type issue5results struct {
	Closed      bool
	Sequence    int
	CloseError  bool // true to return an error on close
	CreateError bool // true to return an error on create
	Created     bool // true if created
}

func (v *issue5resource) test(t *testing.T, expected issue5results) {
	result := v.issue5results
	if result.Closed != expected.Closed {
		if expected.Closed {
			t.Errorf("resource %d not close", v.Id)
		} else {
			t.Errorf("resource %d closed when not expected", v.Id)
		}
	}

	if result.Sequence != expected.Sequence {
		t.Errorf("resource %d expected %d got %d", v.Id, expected.Sequence, result.Sequence)
	}

	if result.CloseError != expected.CloseError {
		t.Errorf("resource %d close error expected %v got %v", v.Id, expected.CloseError, result.CloseError)
	}

	if result.Created != expected.Created {
		t.Errorf("resource %d created expected %v got %v", v.Id, expected.Created, result.Created)
	}

	if result.CreateError != expected.CreateError {
		t.Errorf("resource %d create error expected %v got %v", v.Id, expected.CreateError, result.CreateError)
	}
}

func (v *issue5resource) Close() error {
	v.issue5results.Closed = true
	v.issue5results.Sequence = *v.Counter
	*v.Counter++
	return nil
}

type issue5createCloser struct {
	resource *issue5resource
}

func (v *issue5createCloser) Create() error {
	if v.resource.CreateError {
		v.resource.issue5results.Created = false
		v.resource.issue5results.CreateError = true
		return fmt.Errorf("resource %d create failed", v.resource.Id)
	}
	v.resource.issue5results.CreateError = false
	v.resource.issue5results.Created = true
	return nil
}

func (v *issue5createCloser) Close() error {
	return v.resource.Close()
}

// Test_issue_5 tests that we close resources in the correct order
func Test_issue_5(t *testing.T) {

	tests := []struct {
		name      string
		script    string
		resources []*issue5resource
		results   []issue5results
	}{
		// Test single resource works
		{
			name:   "single",
			script: `main() { try( res1 ) { } }`,
			resources: []*issue5resource{
				{Id: 1},
			},
			results: []issue5results{
				{
					Closed:   true,
					Sequence: 1,
				},
			},
		},
		// Test both resources get closed
		{
			name:   "two",
			script: `main() { try( res1; res2 ) { } }`,
			resources: []*issue5resource{
				{Id: 1},
				{Id: 2},
			},
			results: []issue5results{
				{
					Closed:   true,
					Sequence: 2, // res1 should be closed AFTER res2
				},
				{
					Closed:   true,
					Sequence: 1, // res2 should be closed first as declared last
				},
			},
		},
		// Test first resource is closed when second fails to create
		{
			name:   "error2",
			script: `main() { try( res1; res2 ) { } }`,
			resources: []*issue5resource{
				{Id: 1},
				{
					Id:           2,
					CreateCloser: true,
					CreateError:  true,
				},
			},
			results: []issue5results{
				{
					Closed:   true,
					Sequence: 1, // res1 should be the only one closed
				},
				{
					Closed:      false, // We should not be closed
					CreateError: true,  // We should have got an error
				},
			},
		},
		// Test 3 resources
		{
			name:   "error3-all-close",
			script: `main() { try( res1; res2; res3 ) { } }`,
			resources: []*issue5resource{
				{Id: 1},
				{Id: 2},
				{Id: 3},
			},
			results: []issue5results{
				{
					Closed:   true,
					Sequence: 3, // res1 should be closed last
				},
				{
					Closed:   true,
					Sequence: 2, // res2 should be closed second
				},
				{
					Closed:   true,
					Sequence: 1, // res3 should be closed first
				},
			},
		},
		// Test first resource should close when second fails, third is not created or closed
		{
			name:   "error3-2-fails",
			script: `main() { try( res1; res2; res3 ) { } }`,
			resources: []*issue5resource{
				{Id: 1},
				{
					Id:           2,
					CreateCloser: true,
					CreateError:  true,
				},
				{Id: 3},
			},
			results: []issue5results{
				{
					Closed:   true,
					Sequence: 1, // res1 should be only one closed
				},
				{
					Closed:      false, // We should not be closed
					CreateError: true,  // We should have got an error
				},
				{
					Closed: false, // We should not be closed
				},
			},
		},
		// Test first two resources are closed when third fails to create
		{
			name:   "error3-3-fails",
			script: `main() { try( res1; res2; res3 ) { } }`,
			resources: []*issue5resource{
				{Id: 1},
				{Id: 2},
				{
					Id:           3,
					CreateCloser: true,
					CreateError:  true,
				},
			},
			results: []issue5results{
				{
					Closed:   true,
					Sequence: 2, // res1 should be closed last
				},
				{
					Closed:   true,
					Sequence: 1, // res2 should be closed first
				},
				{
					Closed:      false, // We should not be closed
					CreateError: true,  // We should have got an error
				},
			},
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

			// Sequence we close under, set to first closure value not 0
			// as we can then use 0 to indicate it never ran
			sequence := 1

			// Add each resource to the global scope
			globals := exec.GlobalScope()
			for _, res := range test.resources {
				res.Counter = &sequence
				name := fmt.Sprintf("res%d", res.Id)
				globals.Declare(name)

				switch {
				case res.CreateCloser:
					globals.Set(name, &issue5createCloser{
						resource: res,
					})
				default:
					globals.Set(name, res)
				}
			}

			// Ignore errors as we test for them
			_ = exec.Run()
			//if err != nil {
			//	t.Fatal(err)
			//	return
			//}

			// Test each resource's state
			for i, res := range test.resources {
				res.test(t, test.results[i])
			}
		})
	}
}
