package tests

import (
	"fmt"
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/parser"
	"reflect"
	"strings"
	"testing"
)

func Test_bit_ops(t *testing.T) {
	tests := []struct {
		expr    string
		args    []any
		want    any
		wantErr string
	}{
		// Test <<
		{expr: "%d << %d", args: []any{1, 1}, want: 2},
		{expr: "%d << %d", args: []any{1, 2}, want: 4},
		{expr: "%d << %d", args: []any{1, 3}, want: 8},
		{expr: "%d << %d", args: []any{2, 1}, want: 4},
		{expr: "%d << %d", args: []any{2, 2}, want: 8},
		{expr: "%d << %d", args: []any{2, 3}, want: 16},
		// Test << with floats, this should be unsupported
		{expr: "%f << %d", args: []any{1.0, 1}, want: 2, wantErr: "operation \"float64 << int\" unsupported"},
		{expr: "%f << %d", args: []any{1.5, 1}, want: 2, wantErr: "operation \"float64 << int\" unsupported"},
		{expr: "%d << %f", args: []any{1, 1.0}, want: 2, wantErr: "operation \"int << float64\" unsupported"},
		{expr: "%d << %f", args: []any{1, 1.5}, want: 2, wantErr: "operation \"int << float64\" unsupported"},
		// Test >>
		{expr: "%d >> %d", args: []any{2, 1}, want: 1},
		{expr: "%d >> %d", args: []any{4, 2}, want: 1},
		{expr: "%d >> %d", args: []any{8, 3}, want: 1},
		{expr: "%d >> %d", args: []any{4, 1}, want: 2},
		{expr: "%d >> %d", args: []any{8, 2}, want: 2},
		{expr: "%d >> %d", args: []any{16, 3}, want: 2},
		// Test >> with floats, this should be unsupported
		{expr: "%f >> %d", args: []any{1.0, 1}, want: 2, wantErr: "operation \"float64 >> int\" unsupported"},
		{expr: "%f >> %d", args: []any{1.5, 1}, want: 2, wantErr: "operation \"float64 >> int\" unsupported"},
		{expr: "%d >> %f", args: []any{1, 1.0}, want: 2, wantErr: "operation \"int >> float64\" unsupported"},
		{expr: "%d >> %f", args: []any{1, 1.5}, want: 2, wantErr: "operation \"int >> float64\" unsupported"},
	}

	for tid, tt := range tests {
		fArgs := append([]any{tid}, tt.args...)
		tn := fmt.Sprintf("%d "+tt.expr, fArgs...)
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
					if tt.wantErr == "" || !strings.Contains(err.Error(), tt.wantErr) {
						t.Error(err)
					}
					return
				} else if tt.wantErr != "" {
					t.Errorf("Expected error %q got none", tt.wantErr)
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
