package executor

import (
	"fmt"
	"github.com/peter-mount/go-script/calculator"
	"github.com/peter-mount/go-script/errors"
	"github.com/peter-mount/go-script/script"
	"reflect"
)

func (e *executor) callFunc(cf *script.CallFunc) error {
	err := e.callFuncImpl(cf)

	// Handle return values
	if ret, ok := err.(*errors.ReturnError); ok {
		e.calculator.Push(ret.Value())
		return nil
	}

	return err
}

func (e *executor) callFuncImpl(cf *script.CallFunc) error {

	// Lookup builtin functions
	libFunc, exists := Lookup(cf.Name)
	if exists {
		return libFunc(e, cf)
	}

	// Lookup local function
	f, exists := e.state.GetFunction(cf.Pos, cf.Name)
	if !exists {
		return fmt.Errorf("%s function %q not defined", cf.Pos, cf.Name)
	}

	args, err := e.ProcessParameters(cf)
	if err != nil {
		return err
	}

	return errors.Error(f.Pos, e.functionImpl(f, args))
}

func (e *executor) ProcessParameters(cf *script.CallFunc) ([]interface{}, error) {

	// Process parameters
	var args []interface{}
	if cf.Parameters != nil {
		for _, p := range cf.Parameters.Args {
			v, ok, err := e.calculator.Calculate(func() error {
				return e.assignment(p.Right)
			})
			if err != nil {
				return nil, errors.Error(p.Pos, err)
			}
			if !ok {
				return nil, errors.Errorf(p.Pos, "No result from argument")
			}
			args = append(args, v)
		}
	}

	return args, nil
}

// functionImpl invokes a function declared within the script.
// Used by callFuncImpl and executor.Run
func (e *executor) functionImpl(f *script.FuncDec, args []interface{}) error {
	// Use NewRootScope so we cannot access variables outside the function
	e.state.NewRootScope()

	// Set the current function to this one preserving the caller
	oldFunc := e.state.SetFunction(f)

	// Restore state once we complete
	defer func() {
		e.state.EndScope()
		e.state.SetFunction(oldFunc)
	}()

	if len(args) != len(f.Parameters) {
		return fmt.Errorf("%s parameter mismatch, expected %d got %d", f.Pos, len(f.Parameters), len(args))
	}

	for i, p := range f.Parameters {
		e.state.Declare(p)
		e.state.Set(p, args[i])
	}

	return errors.Error(f.Pos, e.Statements(f.FunBody))
}

// callReflectFunc invokes a function within go from a script
func (e *executor) callReflectFunc(cf *script.CallFunc, f reflect.Value) (ret interface{}, err error) {
	args, err := e.ProcessParameters(cf)
	if err != nil {
		return nil, err
	}

	return e.CallReflectFuncImpl(cf, f, args)
}

// CallReflectFuncImpl makes a function call via reflection.
// Used by callReflectFunc and tests
func (e *executor) CallReflectFuncImpl(cf *script.CallFunc, f reflect.Value, args []interface{}) (ret interface{}, err error) {
	// Any panics get resolved to errors
	defer func() {
		if err1 := recover(); err1 != nil {
			err = errors.Errorf(cf.Pos, "%v", err1)
		}
	}()

	tf := f.Type()

	argVals, err := e.ArgsToValues(cf, tf, args)
	if err != nil {
		return nil, err
	}

	ret0 := f.Call(argVals)

	ret1, err := e.valuesToRet(tf, ret0)
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
		return ret1, nil
	}
}

func (e *executor) ArgsToValues(cf *script.CallFunc, tf reflect.Type, args []interface{}) (ret []reflect.Value, err error) {
	// Any panics get resolved to errors
	defer func() {
		if err1 := recover(); err1 != nil {
			err = errors.Errorf(cf.Pos, "ArgsToValues %v", err1)
		}
	}()

	// '...' so if the last value in args so if it's a slice expand it
	// Unlike go this is supported for any function call not just variadic
	if cf.Parameters != nil && cf.Parameters.Variadic {
		if len(args) == 0 {
			return nil, errors.Errorf(cf.Pos, "'...' with no arguments")
		}

		lArg := args[len(args)-1]
		lArgV := reflect.ValueOf(lArg)
		if lArgV.Kind() == reflect.Slice {
			lArgC := lArgV.Len()

			// Replace last arg with the slice's content
			args = args[:len(args)-1]
			for i := 0; i < lArgC; i++ {
				args = append(args, lArgV.Index(i).Interface())
			}
		}
	}

	// argC = number of arguments function accepts.
	// However, if it's variadic when we take the last one off as we handle that
	// last one specially due to it being a slice.
	argC := tf.NumIn()
	if tf.IsVariadic() {
		argC = argC - 1
	}

	// Cast every argument excluding the variadic one (if present)
	for argN, argV := range args {
		if argN < argC {
			ret, err = e.castArg(ret, argV, tf.In(argN))
			if err != nil {
				return nil, errors.Error(cf.Parameters.Args[argN].Pos, err)
			}
		}
	}

	if tf.IsVariadic() {
		// Type of variadic parameter, always the last one
		variadicIndex := tf.NumIn() - 1
		// last element is actually a slice so call Elem to get the actual type
		variadicType := tf.In(variadicIndex).Elem()

		// For remaining args convert to the same type as the Variadic
		for i := len(ret); i < len(args); i++ {
			ret, err = e.castArg(ret, args[i], variadicType)
			if err != nil {
				return nil, errors.Error(cf.Parameters.Args[variadicIndex].Pos, err)
			}
		}
	}

	return
}

func (e *executor) castArg(ret []reflect.Value, arg interface{}, as reflect.Type) ([]reflect.Value, error) {
	argV := reflect.ValueOf(arg)
	val, err := calculator.Cast(argV, as)
	if err != nil {
		return nil, err
	}
	return append(ret, val), nil
}

func (e *executor) valuesToRet(tf reflect.Type, retVal []reflect.Value) (ret []interface{}, err error) {

	for i := 0; i < tf.NumOut(); i++ {
		tOut := tf.Out(i)
		rv := retVal[i]

		if tOut.Implements(errorInterface) {
			// if err not nil fail the function
			// otherwise drop the value from the results
			if !rv.IsNil() {
				v := rv.Interface()
				return nil, v.(error)
			}
		} else {
			ret = append(ret, rv.Interface())
		}
	}

	return
}

func (e *executor) returnStatement(ret *script.Return) error {
	if ret.Result == nil {
		return errors.NewReturn(nil)
	}

	v, ok, err := e.calculator.Calculate(func() error {
		return e.Expression(ret.Result)
	})
	if err != nil {
		return errors.Error(ret.Pos, err)
	}
	if !ok {
		return errors.Errorf(ret.Pos, "No result from argument")
	}

	return errors.NewReturn(v)
}

var (
	errorInterface = reflect.TypeOf((*error)(nil)).Elem()
)
