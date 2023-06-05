package executor

import (
	"context"
	"fmt"
	"github.com/peter-mount/go-script/calculator"
	"github.com/peter-mount/go-script/script"
)

func (e *executor) callFunc(ctx context.Context) error {
	cf := script.CallFuncFromContext(ctx)

	libFunc, exists := Lookup(cf.Name)
	if exists {
		return libFunc(e, cf, ctx)
	}

	f, exists := e.state.GetFunction(cf.Name)
	if !exists {
		return fmt.Errorf("%s function %q not defined", cf.Pos, cf.Name)
	}

	// Todo parameters
	return e.function(f)
}

func (e *executor) function(f *script.FuncDec, args ...interface{}) error {
	err := e.functionImpl(f, args)

	// Handle return values
	if ret, ok := err.(*returnError); ok {
		if f.ReturnType != "" && f.ReturnType != "void" {
			v, err := calculator.GetValue(f.ReturnType, ret.Value)
			if err != nil {
				return err
			}
			e.calculator.Push(v)
		}
		return nil
	}

	// Should not happen but capture breaks, so they don't leak out of the function
	if IsBreak(err) {
		return nil
	}

	return err
}

func (e *executor) functionImpl(f *script.FuncDec, args []interface{}) error {
	e.state.NewScope()
	defer e.state.EndScope()

	if len(args) != len(f.Parameters) {
		return fmt.Errorf("%s parameter mismatch", f.Pos)
	}

	for i, p := range f.Parameters {
		if p.Scalar != nil {
			e.state.Declare(p.Scalar.Name)
			e.state.Set(p.Scalar.Name, args[i])
		} else if p.Array != nil {
			e.state.Declare(p.Array.Name)
			e.state.Set(p.Array.Name, args[i])
		}
	}

	body := f.FunBody
	if body.Locals != nil {
		for _, l := range body.Locals {
			if l.ScalarDec != nil {
				e.state.Declare(l.ScalarDec.Name)
			} else if l.ArrayDec != nil {
				e.state.Declare(l.ArrayDec.Name)
			}
		}
	}

	return e.visitor.VisitStatements(body.Statements)
	//return e.calculator.Exec(e.statements, body.Statements.WithContext(e.context))
}
