package executor

import (
	"context"
	"fmt"
	"github.com/peter-mount/go-script/calculator"
	"github.com/peter-mount/go-script/script"
	"reflect"
	"sync"
)

var (
	mutex   sync.Mutex
	library = map[string]Function{}
)

type Function func(e Executor, call *script.CallFunc, ctx context.Context) error

func (f Function) Do(e Executor, call *script.CallFunc, ctx context.Context) error {
	if f == nil {
		return nil
	}
	return f(e, call, ctx)
}

func (f Function) Then(b Function) Function {
	if f == nil {
		return b
	}
	if b == nil {
		return f
	}
	return func(e Executor, call *script.CallFunc, ctx context.Context) error {
		if err := f(e, call, ctx); err != nil {
			return err
		}
		return b(e, call, ctx)
	}
}

func Of(funcs ...Function) Function {
	var r Function
	for _, f := range funcs {
		r = r.Then(f)
	}
	return r
}

// Register a Function against a name.
// This will panic if name has already been registered
func Register(name string, f Function) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := library[name]; exists {
		panic(fmt.Errorf("function %q already registered", name))
	}

	library[name] = f
}

// Lookup a registered Function by name
func Lookup(name string) (Function, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	f, exists := library[name]
	if exists {
		return f, true
	}

	return nil, false
}

// RegisterFloat1 registers a function that accepts a float64 as its argument and returns a float64.
// This is common for mathematical functions
func RegisterFloat1(name string, f func(float64) float64) {
	Register(name, func(e Executor, call *script.CallFunc, ctx context.Context) error {
		return Error(call.Pos, float1(f, e, call, ctx))
	})
}

// RegisterFloat2 registers a function that accepts 2 float64's as arguments and returns a float64.
// This is common for mathematical functions
func RegisterFloat2(name string, f func(float64, float64) float64) {
	Register(name, func(e Executor, call *script.CallFunc, ctx context.Context) error {
		return Error(call.Pos, float2(f, e, call, ctx))
	})
}

func float1(f func(float64) float64, e Executor, call *script.CallFunc, ctx context.Context) error {
	arg, err := Args(e, call, ctx)
	if err != nil {
		return err
	}
	if len(arg) != 1 {
		return fmt.Errorf("expected 1 arg, got %d", len(arg))
	}

	af, err := calculator.GetFloat(arg[0])
	if err != nil {
		return err
	}

	return NewReturn(f(af))
}

func float2(f func(float64, float64) float64, e Executor, call *script.CallFunc, ctx context.Context) error {
	arg, err := Args(e, call, ctx)
	if err != nil {
		return err
	}
	if len(arg) != 2 {
		return fmt.Errorf("expected 2 args, got %d", len(arg))
	}

	a, b, err := calculator.Convert(arg[0], arg[1])
	if err != nil {
		return err
	}

	af, err := calculator.GetFloat(a)
	if err != nil {
		return err
	}

	bf, err := calculator.GetFloat(b)
	if err != nil {
		return err
	}

	return NewReturn(f(af, bf))
}

func FuncDelegate(f any) Function {
	fV := reflect.ValueOf(f)
	fT := fV.Type()
	if fT.Kind() != reflect.Func {
		panic(fmt.Errorf("cannot create function delegate for %T", f))
	}

	return func(e Executor, call *script.CallFunc, ctx context.Context) error {

		// Validate argument count
		argC := 0
		if call.Parameters != nil {
			argC = len(call.Parameters.Args)
		}

		numIn := fT.NumIn()

		if fT.IsVariadic() && argC < (numIn-1) {
			return Errorf(call.Pos, "%n requires at least %d parameters", call.Name, numIn)
		}
		if !fT.IsVariadic() && argC != fT.NumIn() {
			return Errorf(call.Pos, "%n requires %d parameters", call.Name, numIn)
		}

		// Process arguments
		args, err := e.ProcessParameters(call, ctx)
		if err != nil {
			return Error(call.Pos, err)
		}

		// call method
		ret, err := e.CallReflectFuncImpl(call, fV, args)
		if err != nil {
			return Error(call.Pos, err)
		}

		// Result on stack
		e.Calculator().Push(ret)
		return nil
	}
}
