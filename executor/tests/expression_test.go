package tests

import (
	"fmt"
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/parser"
	"reflect"
	"testing"
)

func Test_executor_Addition(t *testing.T) {
	tests := []struct {
		expr    string
		args    []any
		want    any
		wantErr bool
	}{
		// A + B - C - D however early on it was calculating (A + (B - (C - D) ) ) specifically A B C D - - +
		// when it should be (((A + B) - C) - D) or A B + C - D -
		// so ensure that the processing order is correct
		// Fixed 2023-06-18
		{expr: "%v+%v-%v-%v", args: []any{1, 2, 3, 4}, want: 1 + 2 - 3 - 4, wantErr: false},
		{expr: "%v+%v-%v-%v", args: []any{1, 2, 6, 3}, want: 1 + 2 - 6 - 3, wantErr: false},
		// A + B * C - D should be A + (B * C) - D but currently processes as ((A + B)*C)-D)
		// Broken 2023-06-18
		{expr: "%v+%v*%v-%v", args: []any{1, 2, 6, 3}, want: 1 + 2*6 - 3, wantErr: false},
	}

	for _, tt := range tests {
		tn := fmt.Sprintf(tt.expr, tt.args...)
		t.Run(tn,
			func(t *testing.T) {
				src := fmt.Sprintf("main() {result = "+tt.expr+"}", tt.args...)

				// Uncomment for debugging
				//fmt.Println(src)

				s, err := parser.New().ParseString(tn, src)
				if err != nil {
					t.Error(err)
					return
				}

				e, err := executor.New(s)
				if err != nil {
					t.Error(err)
					return
				}

				globals := e.GlobalScope()
				globals.Declare("result")

				err = e.Run()
				if err != nil {
					if !tt.wantErr {
						t.Error(err)
					}
					return
				}

				got, exists := globals.Get("result")
				if !exists {
					t.Error("No result returned")
				} else if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Got %T %v wanted %T %v", got, got, tt.want, tt.want)
				}

			})
	}
}
