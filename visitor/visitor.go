package visitor

import (
	"context"
	"github.com/peter-mount/go-kernel/v2/util/task"
	"github.com/peter-mount/go-script/script"
)

type Visitor interface {
	VisitAddition(s *script.Addition) error
	VisitAssignment(s *script.Assignment) error
	VisitCallFunc(s *script.CallFunc) error
	VisitComparison(s *script.Comparison) error
	VisitEquality(s *script.Equality) error
	VisitExpression(s *script.Expression) error
	VisitFuncDec(s *script.FuncDec) error
	VisitFuncBody(body *script.FuncBody) error
	VisitMultiplication(s *script.Multiplication) error
	VisitParameter(p *script.Parameter) error
	VisitPrimary(s *script.Primary) error
	VisitScript(script *script.Script) error
	VisitStatement(s *script.Statement) error
	VisitStatements(s *script.Statements) error
	VisitUnary(s *script.Unary) error
}

func FromContext(ctx context.Context) Visitor {
	return ctx.Value(contextKey).(Visitor)
}

// visitorCommon shared between visitor & visitorBuilder, so it's
// setup with 1 definition in visitorBuilder.WithContext with
// go handling the copying, so we never miss an entry
type visitorCommon struct {
	addition         task.Task
	arrayDec         task.Task
	arrayParamDec    task.Task
	assignment       task.Task
	callFunc         task.Task
	comparison       task.Task
	equality         task.Task
	expression       task.Task
	expressionNoNest bool
	funcDec          task.Task
	funcBody         task.Task
	multiplication   task.Task
	primary          task.Task
	scalarDec        task.Task
	scalarParamDec   task.Task
	script           task.Task
	statement        task.Task
	statements       task.Task
	statementsNoNest bool
	unary            task.Task
}

type visitor struct {
	visitorCommon
	ctx context.Context
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
