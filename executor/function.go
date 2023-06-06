package executor

import (
	"context"
	"fmt"
	"github.com/peter-mount/go-script/script"
)

var ctr = 0

func (e *executor) callFunc(ctx context.Context) error {
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

	ctr++
	if ctr > 3 {
		panic("boo")
	}

	// Process parameters
	var a []interface{}
	for _, p := range cf.Args {
		v, ok, err := e.calculator.Calculate(e.assignment, p.Right.WithContext(ctx))
		if err != nil {
			return Error(p.Pos, err)
		}
		if !ok {
			return Errorf(p.Pos, "No result from argument")
		}
		a = append(a, v)
	}
	return e.function(f, a...)
}

func (e *executor) function(f *script.FuncDec, args ...interface{}) error {
	err := e.functionImpl(f, args)

	// Handle return values
	if ret, ok := err.(*ReturnError); ok {
		e.calculator.Push(ret.Value())
		return nil
	}

	// Should not happen but capture breaks, so they don't leak out of the function
	if IsBreak(err) {
		return nil
	}

	return Error(f.Pos, err)
}

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

	return e.visitor.VisitStatements(f.FunBody.Statements)
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
