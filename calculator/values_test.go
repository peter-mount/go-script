package calculator

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCast(t *testing.T) {
	tests := []struct {
		from    reflect.Value
		want    reflect.Value
		wantErr bool
	}{
		{reflect.ValueOf(5), reflect.ValueOf(5.0), false},
		{reflect.ValueOf(5), reflect.ValueOf(5), false},
		{reflect.ValueOf(5.5), reflect.ValueOf(5), false},
		{reflect.ValueOf("5"), reflect.ValueOf(5.0), false},
		{reflect.ValueOf("5"), reflect.ValueOf("5"), false},
		{reflect.ValueOf("5"), reflect.ValueOf(5), false},
		{reflect.ValueOf("5.5"), reflect.ValueOf(5.5), false},
		{reflect.ValueOf("5.5"), reflect.ValueOf(5), true},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v %v", tt.from.Kind(), tt.want.Kind()), func(t *testing.T) {
			got, err := Cast(tt.from, tt.want.Type())
			switch {
			case err != nil:
				if !tt.wantErr {
					t.Errorf("Cast() error = %v, wantErr %v", err, tt.wantErr)
				}

			case tt.wantErr:
				t.Errorf("Cast() passed but wanted error")

			case got.IsZero():
				t.Errorf("Cast() got zero want %v", tt.want)

			default:
				if got.Type() != tt.want.Type() {
					t.Errorf("Cast() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestGetFloatRaw(t *testing.T) {
	tests := []struct {
		name string
		src  interface{}
		want float64
		fail bool
	}{
		{name: "float64", src: 1.0, want: 1.0},
		{name: "float32", src: float32(1.0), want: 1.0},
		{name: "int", src: 1, fail: true},
		{name: "string float", src: "1.0", fail: true},
		{name: "string int", src: "1", fail: true},
		{name: "bool", src: true, fail: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := GetFloatRaw(tt.src)
			switch {
			case tt.fail && ok:
				t.Error("expected failure got pass")
			case !tt.fail && !ok:
				t.Errorf("expected ok but got failure")
			case !tt.fail && ok && got != tt.want:
				t.Errorf("GetFloatRaw() got = %v, want %v", got, tt.want)
			default:
				// pass
			}
		})
	}
}
