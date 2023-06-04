package visitor

import (
	"context"
	"github.com/peter-mount/go-kernel/v2/util/task"
	"github.com/peter-mount/go-script/script"
)

type Visitor interface {
	VisitScript(script *script.Script) error
	VisitFuncDec(s *script.FuncDec) error
	VisitFuncBody(body *script.FuncBody) error
	VisitParameter(p *script.Parameter) error
}

func FromContext(ctx context.Context) Visitor {
	return ctx.Value(contextKey).(Visitor)
}

type visitor struct {
	ctx            context.Context
	script         task.Task
	funcDec        task.Task
	funcBody       task.Task
	arrayParamDec  task.Task
	scalarParamDec task.Task
	arrayDec       task.Task
	scalarDec      task.Task
}

func (v *visitor) visit(p func(context.Context) context.Context, f func() error) error {
	oldCtx := v.ctx
	newCtx := p(v.ctx)
	v.ctx = newCtx
	defer func() {
		v.ctx = oldCtx
	}()
	return f()
}

func (v *visitor) visitTask(p func(context.Context) context.Context, t task.Task) error {
	if t == nil {
		return nil
	}

	return v.visit(p, func() error {
		return t.Do(v.ctx)
	})
}

func (v *visitor) VisitScript(s *script.Script) error {
	return v.visit(s.WithContext, func() error {
		if err := v.script.Do(v.ctx); err != nil {
			return err
		}

		for _, e := range s.TopDec {
			if err := v.visitTopDec(e); err != nil {
				return err
			}
		}

		return nil
	})
}

// visitTopDec handles TopDec and VarDec delegating to funcDec, arrayDec or scalarDec
func (v *visitor) visitTopDec(topDec *script.TopDec) error {
	switch {
	case topDec.FunDec != nil:
		return v.VisitFuncDec(topDec.FunDec)
	case topDec.VarDec != nil && topDec.VarDec.ArrayDec != nil:
		return v.visitTask(topDec.VarDec.ArrayDec.WithContext, v.arrayDec)
	case topDec.VarDec != nil && topDec.VarDec.ScalarDec != nil:
		return v.visitTask(topDec.VarDec.ScalarDec.WithContext, v.scalarDec)
	default:
	}
	return nil
}

func (v *visitor) VisitFuncDec(s *script.FuncDec) error {
	return v.visit(s.WithContext, func() error {
		if err := v.funcDec.Do(v.ctx); err != nil {
			return err
		}

		for _, param := range s.Parameters {
			if err := v.VisitParameter(param); err != nil {
				return err
			}
		}

		return nil
	})
}

func (v *visitor) VisitFuncBody(s *script.FuncBody) error {
	return v.visitTask(s.WithContext, v.funcBody)
}

func (v *visitor) VisitParameter(p *script.Parameter) error {
	switch {
	case p.Scalar != nil:
		return v.visitTask(p.Scalar.WithContext, v.scalarParamDec)
	case p.Array != nil:
		return v.visitTask(p.Array.WithContext, v.arrayParamDec)
	default:
		return nil
	}
}
