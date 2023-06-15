package executor

import (
	"context"
	"github.com/peter-mount/go-script/calculator"
	"github.com/peter-mount/go-script/script"
	"reflect"
)

// getReference resolves a reference of a value.
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
func (e *executor) getReference(op *script.Primary, v interface{}, ctx context.Context) (err error) {
	// These are not valid at this point.
	switch {

	case op.Ident != "":
		v, err = e.resolveReference(op, op.Ident, v)
		if err == nil {
			// Handle arrays
			v, err = e.resolveArray(op, v)
		}

	// method reference against v not declared functions
	case op.CallFunc != nil:
		v, err = e.resolveFunction(op.CallFunc, v, ctx)

	// Nonsensical to be part of a reference
	case op.SubExpression != nil:
		return Errorf(op.Pos, "invalid reference")

	// Default say unimplemented as we may allow these in the future?
	default:
		return Errorf(op.Pos, "not supported yet")
	}

	if err != nil {
		return Error(op.Pos, err)
	}

	// recurse as we have a pointer to the next field
	if op.Pointer != nil {
		return e.getReference(op.Pointer, v, ctx)
	}

	e.calculator.Push(v)
	return nil
}

func (e *executor) resolveReference(op *script.Primary, name string, v interface{}) (ret interface{}, err error) {
	// Any panics get resolved to errors
	defer func() {
		if err1 := recover(); err1 != nil {
			err = Errorf(op.Pos, "%v", err1)
		}
	}()

	tv := reflect.ValueOf(v)
	ti := reflect.Indirect(tv)

	switch ti.Kind() {
	case reflect.Struct:
		tf := ti.FieldByName(name)
		if tf.IsValid() {
			ret = tf.Interface()
			return
		}

	case reflect.Map:
		me := ti.MapIndex(reflect.ValueOf(name))
		if me.IsValid() {
			ret = me.Interface()
			return
		}

	case reflect.Array, reflect.Slice, reflect.String:
		var idx int
		idx, err = calculator.GetInt(name)
		if err != nil {
			return nil, Error(op.Pos, err)
		}

		if idx < 0 || idx >= ti.Len() {
			return nil, Errorf(op.Pos, "Index out of bounds %d", idx)
		}

		ret = ti.Index(idx).Interface()
		return
	}

	return nil, Errorf(op.Pos, "%T has no field %q", v, name)
}

func (e *executor) resolveArray(op *script.Primary, v interface{}) (interface{}, error) {
	// Nothing to do
	if len(op.ArrayIndex) == 0 {
		return v, nil
	}

	// Run through each dimension and set v to the result of each lookup
	for _, dimension := range op.ArrayIndex {
		err := e.expression(dimension.WithContext(e.context))
		if err != nil {
			return nil, Error(dimension.Pos, err)
		}

		index, err := e.calculator.Pop()
		if err != nil {
			return nil, Error(dimension.Pos, err)
		}

		v, err = e.resolveArrayIndex(index, v, dimension)
		if err != nil {
			return nil, Error(dimension.Pos, err)
		}
	}

	return v, nil
}

func (e *executor) resolveArrayIndex(index, v interface{}, dimension *script.Expression) (ret interface{}, err error) {
	// Any panics get resolved to errors
	defer func() {
		if err1 := recover(); err1 != nil {
			err = Errorf(dimension.Pos, "%v", err1)
		}
	}()

	tv := reflect.ValueOf(v)
	ti := reflect.Indirect(tv)

	switch ti.Kind() {

	case reflect.Array, reflect.Slice, reflect.String:
		var idx int
		idx, err = calculator.GetInt(index)
		if err != nil {
			return nil, Error(dimension.Pos, err)
		}

		if idx < 0 || idx >= ti.Len() {
			return nil, Errorf(dimension.Pos, "Index out of bounds %d", idx)
		}

		ret = ti.Index(idx).Interface()

	case reflect.Struct:
		var name string
		name, err = calculator.GetString(index)
		if err != nil {
			return nil, Error(dimension.Pos, err)
		}

		tf := ti.FieldByName(name)
		if tf.IsValid() {
			ret = tf.Interface()
		} else {
			return nil, Errorf(dimension.Pos, "%T has no field %q", v, name)
		}

	case reflect.Map:
		me := ti.MapIndex(reflect.ValueOf(index))
		if me.IsValid() {
			ret = me.Interface()
		} else {
			return nil, Errorf(dimension.Pos, "%T has no field %v", v, v)
		}

	default:
		ret = v
	}

	return
}

func (e *executor) resolveFunction(op *script.CallFunc, v interface{}, ctx context.Context) (ret interface{}, err error) {

	tv := reflect.ValueOf(v)
	// Loop so we check tv, then if that doesn't match *tv
	// This allows for "func(m *type) name()" and "func(m type) name()"
	for _, ti := range []reflect.Value{tv, reflect.Indirect(tv)} {
		switch ti.Kind() {
		case reflect.Struct:
			tf := tv.MethodByName(op.Name)
			if tf.IsValid() {
				ret, err = e.callReflectFunc(op, tf, ctx)
				return
			}
		}
	}

	return nil, Errorf(op.Pos, "%T has no function %q", v, op.Name)
}
