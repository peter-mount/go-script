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
	VisitStatements(s *script.Statements) error
	VisitStatement(s *script.Statement) error
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
	statements     task.Task
	statement      task.Task
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
