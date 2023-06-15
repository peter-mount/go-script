package executor

import (
	"context"
	"fmt"
	"github.com/peter-mount/go-script/calculator"
	"github.com/peter-mount/go-script/script"
	"reflect"
)

func (e *executor) callFunc(ctx context.Context) error {
	err := e.callFuncImpl(ctx)
	// Handle return values
	if ret, ok := err.(*ReturnError); ok {
		e.calculator.Push(ret.Value())
		return nil
	}

	// Should not happen but capture breaks, so they don't leak out of the function
	if IsBreak(err) {
		return nil
	}

	return err
}

func (e *executor) callFuncImpl(ctx context.Context) error {
	cf := script.CallFuncFromContext(ctx)

	// Lookup builtin functions
	libFunc, exists := Lookup(cf.Name)
	if exists {
		return libFunc(e, cf, ctx)
	}

	// Lookup local function
	f, exists := e.state.GetFunction(cf.Name)
	if !exists {
		return fmt.Errorf("%s function %q not defined", cf.Pos, cf.Name)
	}

	args, err := e.processParameters(cf, ctx)
	if err != nil {
		return err
	}

	return Error(f.Pos, e.functionImpl(f, args))
}

func (e *executor) processParameters(cf *script.CallFunc, ctx context.Context) ([]interface{}, error) {

	// Process parameters
	var args []interface{}
	for _, p := range cf.Args {
		v, ok, err := e.calculator.Calculate(e.assignment, p.Right.WithContext(ctx))
		if err != nil {
			return nil, Error(p.Pos, err)
		}
		if !ok {
			return nil, Errorf(p.Pos, "No result from argument")
		}
		args = append(args, v)
	}

	return args, nil
}

// functionImpl invokes the function.
// Used by callFuncImpl and executor.Run
func (e *executor) functionImpl(f *script.FuncDec, args []interface{}) error {
	// Use NewRootScope so we cannot access variables outside the function
	e.state.NewRootScope()
	defer e.state.EndScope()

	if len(args) != len(f.Parameters) {
		return fmt.Errorf("%s parameter mismatch, expected %d got %d", f.Pos, len(f.Parameters), len(args))
	}

	for i, p := range f.Parameters {
		e.state.Declare(p.Ident)
		e.state.Set(p.Ident, args[i])
	}

	return Error(f.Pos, e.visitor.VisitStatements(f.FunBody.Statements))
}

func (e *executor) callReflectFunc(cf *script.CallFunc, f reflect.Value, ctx context.Context) (ret interface{}, err error) {
	// Any panics get resolved to errors
	defer func() {
		if err1 := recover(); err1 != nil {
			err = Errorf(cf.Pos, "%v", err1)
		}
	}()

	args, err := e.processParameters(cf, ctx)
	if err != nil {
		return nil, err
	}

	tf := f.Type()

	argVals, err := e.argsToValues(cf, tf, args)
	if err != nil {
		return nil, err
	}

	retVal := f.Call(argVals)

	ret1, err := e.valuesToRet(cf, tf, retVal)
	if err != nil {
		return nil, err
	}

	// Work out what to return
	switch len(ret1) {
	case 0:
		return nil, nil

	case 1:
		return ret1[0], nil

	default:
		return ret, nil
	}
}

func (e *executor) argsToValues(cf *script.CallFunc, tf reflect.Type, args []interface{}) (ret []reflect.Value, err error) {
	// Any panics get resolved to errors
	defer func() {
		if err1 := recover(); err1 != nil {
			err = Errorf(cf.Pos, "%v", err1)
		}
	}()

	for argN, argV := range args {
		val := reflect.ValueOf(argV)

		if argN < tf.NumIn() {
			val, err = calculator.Cast(val, tf.In(argN))
			if err != nil {
				return nil, Error(cf.Args[argN].Pos, err)
			}
		}

		ret = append(ret, val)
	}
	return
}

func (e *executor) valuesToRet(cf *script.CallFunc, tf reflect.Type, retVal []reflect.Value) (ret []interface{}, err error) {

	for i := 0; i < tf.NumOut(); i++ {
		tOut := tf.Out(i)
		tk := tOut.Kind()

		rv := retVal[i]

		if tOut.Implements(errorInterface) {
			// if err not nil fail the function
			// otherwise drop the value from the results
			if !rv.IsNil() {
				v := rv.Interface()
				return nil, v.(error)
			}
		} else {
			var v interface{}

			switch tk {

			case reflect.Float64, reflect.Float32:
				v = rv.Float()

			case reflect.Int, reflect.Int64,
				reflect.Int8, reflect.Int16, reflect.Int32:
				v = rv.Int()

			case reflect.Uint, reflect.Uint64,
				reflect.Uint8, reflect.Uint16, reflect.Uint32:
				v = rv.Int()

			case reflect.Array, reflect.Map, reflect.Struct:
				v = rv.Interface()

			case reflect.String:
				v = rv.String()

			case reflect.Bool:
				v = rv.Bool()

			default:
				v = rv.Interface()
			}

			ret = append(ret, v)
		}
	}

	return
}

func (e *executor) returnStatement(ctx context.Context) error {
	ret := script.ReturnFromContext(ctx)

	if ret.Result == nil {
		return NewReturn(nil)
	}

	v, ok, err := e.calculator.Calculate(e.expression, ret.Result.WithContext(ctx))
	if err != nil {
		return Error(ret.Pos, err)
	}
	if !ok {
		return Errorf(ret.Pos, "No result from argument")
	}

	return NewReturn(v)
}

var (
	errorInterface = reflect.TypeOf((*error)(nil)).Elem()
)
