package state

import (
	"reflect"
	"strings"
)

func (s *state) Get(n string) (interface{}, bool) {
	// Plain get of a variable
	if !strings.Contains(n, ".") {
		return s.variables.Get(n)
	}

	// Split . path & get first entry as a variable
	na := strings.Split(n, ".")
	v, exists := s.variables.Get(na[0])
	if !exists {
		return nil, false
	}

	// Now resolve each "field"
	for _, ne := range na[1:] {
		v, exists = resolveField(v, ne)
		if !exists {
			return nil, false
		}
	}

	return v, true
}

func resolveField(v interface{}, n string) (interface{}, bool) {
	if v == nil {
		return nil, false
	}

	tv := reflect.ValueOf(v)
	ti := reflect.Indirect(tv)

	switch ti.Kind() {
	case reflect.Struct:
		tf := ti.FieldByName(n)
		if tf.IsValid() {
			return tf.Interface(), true
		}
		// todo if not found then method?
		//ti.MethodByName(n)
		return nil, false

	case reflect.Map:
		me := ti.MapIndex(reflect.ValueOf(n))
		if me.IsValid() {
			return me.Interface(), true
		}
		return nil, false

	default:
		return ti.Interface(), false
	}
}
