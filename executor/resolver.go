package executor

import (
	"github.com/peter-mount/go-script/script"
	"reflect"
)

// resolveReference resolves a reference of a value.
//
// This handles pointers, e.g. var.field.subfield.value where var, field, subfield etc.
// can be a map or struct.
//
// Note: For structs the pointer names must start with an upper case character due to
// go only allowing those field names to be public.
//
// op the Primary defining the reference
//
// v the value this Primary is referencing.
func (e *executor) resolveReference(op *script.Primary, v interface{}) error {
	// These are not valid at this point.
	switch {

	case op.Ident != "":
		nv, err := e.resolveRefName(op, op.Ident, v)
		if err != nil {
			return Error(op.Pos, err)
		}

		switch {

		// Move to next pointer in the chain
		case op.Pointer != nil:
			return e.resolveReference(op.Pointer, nv)

		case op.ArrayIndex != nil:
			return Errorf(op.Pos, "arrays unsupported")

		// Push result to stack and finish
		default:
			e.calculator.Push(nv)
			return nil
		}

	// Nonsensical to be part of a reference
	case op.SubExpression != nil:
		return Errorf(op.Pos, "invalid reference")

	// method reference against v not declared functions
	case op.CallFunc != nil:
		return Errorf(op.Pos, "not supported yet")

	// Default say unimplemented as we may allow these in the future?
	default:
		return Errorf(op.Pos, "not supported yet")
	}
}

func (e *executor) resolveRefName(op *script.Primary, name string, v interface{}) (interface{}, error) {

	tv := reflect.ValueOf(v)
	ti := reflect.Indirect(tv)

	switch ti.Kind() {
	case reflect.Struct:
		tf := ti.FieldByName(name)
		if tf.IsValid() {
			return tf.Interface(), nil
		}

	case reflect.Map:
		me := ti.MapIndex(reflect.ValueOf(name))
		if me.IsValid() {
			return me.Interface(), nil
		}

	}

	return nil, Errorf(op.Pos, "%T has no field %q", v, name)
}
