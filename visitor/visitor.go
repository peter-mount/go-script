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
	VisitFor(s *script.For) error
	VisitForRange(s *script.ForRange) error
	VisitFuncDec(s *script.FuncDec) error
	VisitLogic(s *script.Logic) error
	VisitIf(s *script.If) error
	VisitMultiplication(s *script.Multiplication) error
	VisitPrimary(s *script.Primary) error
	VisitReturn(s *script.Return) error
	VisitScript(script *script.Script) error
	VisitStatement(s *script.Statement) error
	VisitStatements(s *script.Statements) error
	VisitTry(s *script.Try) error
	VisitUnary(s *script.Unary) error
	VisitWhile(s *script.While) error
}

func FromContext(ctx context.Context) Visitor {
	return ctx.Value(contextKey).(Visitor)
}

// visitorCommon shared between visitor & visitorBuilder, so it's
// setup with 1 definition in visitorBuilder.WithContext with
// go handling the copying, so we never miss an entry
type visitorCommon struct {
	addition         task.Task
	assignment       task.Task
	callFunc         task.Task
	comparison       task.Task
	equality         task.Task
	expression       task.Task
	expressionNoNest bool
	forRange         task.Task
	forStatement     task.Task
	funcDec          task.Task
	funcBody         task.Task
	ifStatement      task.Task
	logic            task.Task
	multiplication   task.Task
	primary          task.Task
	returnStatement  task.Task
	script           task.Task
	statement        task.Task
	statements       task.Task
	statementsNoNest bool
	try              task.Task
	unary            task.Task
	whileStatement   task.Task
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

		for _, e := range s.FunDec {
			if err := v.VisitFuncDec(e); err != nil {
				return err
			}
		}

		return nil
	})
}
