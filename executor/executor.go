package executor

import (
	"context"
	"fmt"
	"github.com/peter-mount/go-script/calculator"
	"github.com/peter-mount/go-script/script"
	"github.com/peter-mount/go-script/state"
	"github.com/peter-mount/go-script/visitor"
)

type Executor interface {
	Run() error
}

type executor struct {
	script     *script.Script
	state      state.State
	calculator calculator.Calculator
	visitor    visitor.Visitor
	context    context.Context
}

func New(s *script.Script) (Executor, error) {
	execState, err := state.New(s)
	if err != nil {
		return nil, err
	}
	e := &executor{
		script:     s,
		state:      execState,
		calculator: calculator.New(),
	}

	e.context = execState.WithContext(context.Background())

	e.visitor = visitor.New().
		Addition(e.addition).
		Assignment(e.assignment).
		CallFunc(e.callFunc).
		Comparison(e.comparison).
		Equality(e.equality).
		Expression(e.expression).
		ExpressionNoNest().
		Multiplication(e.multiplication).
		Primary(e.primary).
		Statement(e.statement).
		Statements(e.statements).
		StatementsNoNest().
		Unary(e.unary).
		WithContext(e.context)

	return e, nil
}

func (e *executor) Run() error {
	main, hasMain := e.state.GetFunction("main")
	if !hasMain {
		return fmt.Errorf("%s main() function not defined", e.script.Pos)
	}

	err := e.function(main, e.context)

	// Pass err unless it's return or break.
	// break should happen lower down but this catches it, so it doesn't
	// exit the function call
	if err != nil && !(IsReturn(err) || IsBreak(err)) {
		return err
	}

	return nil
}

func (e *executor) callFunc(ctx context.Context) error {
	cf := script.CallFuncFromContext(ctx)
	f, exists := e.state.GetFunction(cf.Name)
	if !exists {
		return fmt.Errorf("%s function %q not defined", cf.Pos, cf.Name)
	}

	// Todo parameters
	return e.function(f, ctx)
}

func (e *executor) function(f *script.FuncDec, ctx context.Context, args ...interface{}) error {
	fmt.Printf("%s exec %q\n", f.Pos, f.Name)

	err := e.functionImpl(f, ctx, args)
	fmt.Printf("%s exec %v\n", f.Pos, err)

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

func (e *executor) functionImpl(f *script.FuncDec, ctx context.Context, args []interface{}) error {
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
