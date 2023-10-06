package tests

import (
	"fmt"
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/parser"
	"reflect"
	"strings"
	"testing"
)

func Test_operations(t *testing.T) {
	tests := []struct {
		expr    string
		args    []any
		want    any
		wantErr string
	}{
		// ==================================================
		// Mono operations, e.g. op val
		// ==================================================
		// Test ! not
		{expr: "!%v", args: []any{true}, want: false},
		{expr: "!%v", args: []any{false}, want: true},
		{expr: "!%d", args: []any{0}, want: true},
		{expr: "!%d", args: []any{1}, want: false},
		{expr: "!%s", args: []any{`"true"`}, wantErr: "operation \"! string\" unsupported"},
		{expr: "!%s", args: []any{`"t"`}, wantErr: "operation \"! string\" unsupported"},
		{expr: "!%s", args: []any{`"false"`}, wantErr: "operation \"! string\" unsupported"},
		{expr: "!%s", args: []any{`"f"`}, wantErr: "operation \"! string\" unsupported"},
		// Test - negation
		{expr: "-%v", args: []any{true}, want: false},
		{expr: "-%v", args: []any{false}, want: true},
		{expr: "-%d", args: []any{0}, want: 0},
		{expr: "-%d", args: []any{1}, want: -1},
		{expr: "-(%d)", args: []any{-10}, want: 10},
		{expr: "-(%d)", args: []any{10}, want: -10},
		{expr: "-%s", args: []any{`"true"`}, wantErr: "operation \"- string\" unsupported"},
		{expr: "-%s", args: []any{`"t"`}, wantErr: "operation \"- string\" unsupported"},
		{expr: "-%s", args: []any{`"false"`}, wantErr: "operation \"- string\" unsupported"},
		{expr: "-%s", args: []any{`"f"`}, wantErr: "operation \"- string\" unsupported"},
		// ==================================================
		// Bi operations, e.g. val op val
		// ==================================================
		// Test ==
		{expr: "%d == %d", args: []any{1, 1}, want: true},
		{expr: "%d == %f", args: []any{1, 1.0}, want: true},
		{expr: "%f == %d", args: []any{1.0, 1}, want: true},
		{expr: "%f == %f", args: []any{1.0, 1.0}, want: true},
		{expr: "%d == %d", args: []any{1, 2}, want: false},
		{expr: "%d == %f", args: []any{1, -10.0}, want: false},
		{expr: "%f == %d", args: []any{1.0, 15}, want: false},
		{expr: "%f == %f", args: []any{1.0, -999.99}, want: false},
		// --------------------------------------------------

		// Test !=
		{expr: "%d != %d", args: []any{1, 1}, want: false},
		{expr: "%d != %f", args: []any{1, 1.0}, want: false},
		{expr: "%f != %d", args: []any{1.0, 1}, want: false},
		{expr: "%f != %f", args: []any{1.0, 1.0}, want: false},
		{expr: "%d != %d", args: []any{1, 2}, want: true},
		{expr: "%d != %f", args: []any{1, -10.0}, want: true},
		{expr: "%f != %d", args: []any{1.0, 15}, want: true},
		{expr: "%f != %f", args: []any{1.0, -999.99}, want: true},
		// --------------------------------------------------

		// Test <
		{expr: "%d < %d", args: []any{1, 1}, want: false},
		{expr: "%d < %f", args: []any{1, 1.0}, want: false},
		{expr: "%f < %d", args: []any{1.0, 1}, want: false},
		{expr: "%f < %f", args: []any{1.0, 1.0}, want: false},
		{expr: "%d < %d", args: []any{1, 2}, want: true},
		{expr: "%d < %d", args: []any{2, 1}, want: false},
		{expr: "%f < %d", args: []any{-10.0, 1}, want: true},
		{expr: "%d < %f", args: []any{1, -10.0}, want: false},
		{expr: "%f < %d", args: []any{1.0, 15}, want: true},
		{expr: "%d < %f", args: []any{15, 1.0}, want: false},
		{expr: "%f < %f", args: []any{-999.99, 1.0}, want: true},
		{expr: "%f < %f", args: []any{1.0, -999.99}, want: false},
		// --------------------------------------------------

		// Test <=
		{expr: "%d <= %d", args: []any{1, 1}, want: true},
		{expr: "%d <= %f", args: []any{1, 1.0}, want: true},
		{expr: "%f <= %d", args: []any{1.0, 1}, want: true},
		{expr: "%f <= %f", args: []any{1.0, 1.0}, want: true},
		{expr: "%d <= %d", args: []any{1, 2}, want: true},
		{expr: "%d <= %d", args: []any{2, 1}, want: false},
		{expr: "%f <= %d", args: []any{-10.0, 1}, want: true},
		{expr: "%d <= %f", args: []any{1, -10.0}, want: false},
		{expr: "%f <= %d", args: []any{1.0, 15}, want: true},
		{expr: "%d <= %f", args: []any{15, 1.0}, want: false},
		{expr: "%f <= %f", args: []any{-999.99, 1.0}, want: true},
		{expr: "%f <= %f", args: []any{1.0, -999.99}, want: false},
		// --------------------------------------------------

		// Test >
		{expr: "%d > %d", args: []any{1, 1}, want: false},
		{expr: "%d > %f", args: []any{1, 1.0}, want: false},
		{expr: "%f > %d", args: []any{1.0, 1}, want: false},
		{expr: "%f > %f", args: []any{1.0, 1.0}, want: false},
		{expr: "%d > %d", args: []any{1, 2}, want: false},
		{expr: "%d > %d", args: []any{2, 1}, want: true},
		{expr: "%f > %d", args: []any{-10.0, 1}, want: false},
		{expr: "%d > %f", args: []any{1, -10.0}, want: true},
		{expr: "%f > %d", args: []any{1.0, 15}, want: false},
		{expr: "%d > %f", args: []any{15, 1.0}, want: true},
		{expr: "%f > %f", args: []any{-999.99, 1.0}, want: false},
		{expr: "%f > %f", args: []any{1.0, -999.99}, want: true},
		// --------------------------------------------------

		// Test >=
		{expr: "%d >= %d", args: []any{1, 1}, want: true},
		{expr: "%d >= %f", args: []any{1, 1.0}, want: true},
		{expr: "%f >= %d", args: []any{1.0, 1}, want: true},
		{expr: "%f >= %f", args: []any{1.0, 1.0}, want: true},
		{expr: "%d >= %d", args: []any{1, 2}, want: false},
		{expr: "%d >= %d", args: []any{2, 1}, want: true},
		{expr: "%f >= %d", args: []any{-10.0, 1}, want: false},
		{expr: "%d >= %f", args: []any{1, -10.0}, want: true},
		{expr: "%f >= %d", args: []any{1.0, 15}, want: false},
		{expr: "%d >= %f", args: []any{15, 1.0}, want: true},
		{expr: "%f >= %f", args: []any{-999.99, 1.0}, want: false},
		{expr: "%f >= %f", args: []any{1.0, -999.99}, want: true},
		// --------------------------------------------------

		// Test + addition
		{expr: "%d + %d", args: []any{1, 1}, want: 2},
		{expr: "%d + %f", args: []any{1, 2.0}, want: 3.0},
		{expr: "%f + %d", args: []any{1.0, 3}, want: 4.0},
		{expr: "%d + %d", args: []any{2, -1}, want: 1},
		{expr: "%d + %d", args: []any{-2, 2}, want: 0},
		{expr: "%d + %d", args: []any{-2, 3}, want: 1},
		// Test with and without whitespace as this then tests against ++ operator.
		// This should work as we are using a number and not an ident
		{expr: "%d + %d", args: []any{2, -1}, want: 1},
		{expr: "%d+%d", args: []any{2, -1}, want: 1},
		// --------------------------------------------------

		// Test - subtraction
		{expr: "%d - %d", args: []any{1, 1}, want: 0},
		{expr: "%f - %d", args: []any{1.0, 2}, want: -1.0},
		{expr: "%d - %f", args: []any{1, 3.0}, want: -2.0},
		{expr: "%d - %d", args: []any{-2, 2}, want: -4},
		{expr: "%d - %d", args: []any{-2, 3}, want: -5},
		// Test with and without whitespace as this then tests against -- operator.
		// This should work as we are using a number and not an ident
		{expr: "%d - %d", args: []any{2, -1}, want: 3},
		{expr: "%d-%d", args: []any{2, -1}, want: 3},
		// --------------------------------------------------

		// Test * multiplication
		{expr: "%d * %d", args: []any{1, 1}, want: 1},
		{expr: "%f * %d", args: []any{1.0, 2}, want: 2.0},
		{expr: "%d * %f", args: []any{1, 3.0}, want: 3.0},
		{expr: "%d * %d", args: []any{-2, 2}, want: -4},
		{expr: "%d * %d", args: []any{-2, 3}, want: -6},
		// --------------------------------------------------

		// Test / division
		{expr: "%d / %d", args: []any{1, 1}, want: 1},
		{expr: "%f / %d", args: []any{1.0, 2}, want: 0.5},
		{expr: "%d / %f", args: []any{1, 3.0}, want: 1 / 3.0},
		{expr: "%d / %d", args: []any{-2, 2}, want: -1},
		{expr: "%d / %d", args: []any{-2, 3}, want: -2 / 3},
		// --------------------------------------------------

		// Test && logical and
		{expr: "%v && %v", args: []any{true, true}, want: true},
		{expr: "%v && %v", args: []any{true, false}, want: false},
		{expr: "%v && %v", args: []any{false, true}, want: false},
		{expr: "%v && %v", args: []any{false, false}, want: false},
		// --------------------------------------------------

		// Test || logical or
		{expr: "%v || %v", args: []any{true, true}, want: true},
		{expr: "%v || %v", args: []any{true, false}, want: true},
		{expr: "%v || %v", args: []any{false, true}, want: true},
		{expr: "%v || %v", args: []any{false, false}, want: false},
		// --------------------------------------------------

		// Test % modulus
		{expr: "%d %% %d", args: []any{1, 1}, want: 0},
		{expr: "%d %% %d", args: []any{1, 2}, want: 1},
		{expr: "%d %% %d", args: []any{6, 5}, want: 1},
		{expr: "%f %% %d", args: []any{1.0, 2}, want: 1.0},
		{expr: "%d %% %f", args: []any{1, 3.0}, want: 1.0},
		{expr: "%d %% %d", args: []any{-2, 2}, want: 0},
		{expr: "%d %% %d", args: []any{-2, 3}, want: -2},
		// --------------------------------------------------

		// Test <<
		{expr: "%d << %d", args: []any{1, 1}, want: 1 << 1},
		{expr: "%d << %d", args: []any{1, 2}, want: 1 << 2},
		{expr: "%d << %d", args: []any{1, 3}, want: 1 << 3},
		{expr: "%d << %d", args: []any{2, 1}, want: 2 << 1},
		{expr: "%d << %d", args: []any{2, 2}, want: 2 << 2},
		{expr: "%d << %d", args: []any{2, 3}, want: 2 << 3},
		// --------------------------------------------------

		// Test << with floats, this should be unsupported
		{expr: "%f << %d", args: []any{1.0, 1}, wantErr: "operation \"float64 << int\" unsupported"},
		{expr: "%f << %d", args: []any{1.5, 1}, wantErr: "operation \"float64 << int\" unsupported"},
		{expr: "%d << %f", args: []any{1, 1.0}, wantErr: "operation \"int << float64\" unsupported"},
		{expr: "%d << %f", args: []any{1, 1.5}, wantErr: "operation \"int << float64\" unsupported"},
		// --------------------------------------------------

		// Test >>
		{expr: "%d >> %d", args: []any{2, 1}, want: 2 >> 1},
		{expr: "%d >> %d", args: []any{4, 2}, want: 4 >> 2},
		{expr: "%d >> %d", args: []any{8, 3}, want: 8 >> 3},
		{expr: "%d >> %d", args: []any{4, 1}, want: 4 >> 1},
		{expr: "%d >> %d", args: []any{8, 2}, want: 8 >> 2},
		{expr: "%d >> %d", args: []any{16, 3}, want: 16 >> 3},
		// --------------------------------------------------

		// Test >> with floats, this should be unsupported
		{expr: "%f >> %d", args: []any{1.0, 1}, wantErr: "operation \"float64 >> int\" unsupported"},
		{expr: "%f >> %d", args: []any{1.5, 1}, wantErr: "operation \"float64 >> int\" unsupported"},
		{expr: "%d >> %f", args: []any{1, 1.0}, wantErr: "operation \"int >> float64\" unsupported"},
		{expr: "%d >> %f", args: []any{1, 1.5}, wantErr: "operation \"int >> float64\" unsupported"},
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
