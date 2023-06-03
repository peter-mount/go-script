package script

import (
	"context"
	"github.com/peter-mount/go-kernel/v2/util/task"
)

type Visitor interface {
	VisitCall(call *Call) error
	VisitPrint(print *Print) error
	VisitRemark(remark *Remark) error
	VisitScript(script *Script) error
	VisitStatement(statement *Statement) error
}

type Visitable interface {
	Accept(v Visitor) error
}

type visitor struct {
	ctx       context.Context
	call      task.Task
	print     task.Task
	remark    task.Task
	script    task.Task
	statement task.Task
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

func (v *visitor) VisitCall(call *Call) error {
	return v.visitTask(call.WithContext, v.call)
}

func (v *visitor) VisitPrint(print *Print) error {
	return v.visitTask(print.WithContext, v.print)
}

func (v *visitor) VisitRemark(remark *Remark) error {
	return v.visitTask(remark.WithContext, v.remark)
}

func (v *visitor) VisitScript(s *Script) error {
	return v.visit(s.WithContext, func() error {
		if err := v.script.Do(v.ctx); err != nil {
			return err
		}

		// FIXME use PC not sequential!
		for _, statement := range s.Statements {
			if err := statement.Accept(v); err != nil {
				return err
			}
		}
		return nil
	})
}

func (v *visitor) VisitStatement(s *Statement) error {
	return v.visit(s.WithContext, func() error {
		if err := v.statement.Do(v.ctx); err != nil {
			return err
		}

		switch {
		case s.Call != nil:
			return s.Call.Accept(v)
		case s.Print != nil:
			return s.Print.Accept(v)
		case s.Remark != nil:
			return s.Remark.Accept(v)
		}

		// TODO raise error here?
		return nil
	})
}

type VisitorBuilder interface {
	Call(task2 task.Task) VisitorBuilder
	Print(t task.Task) VisitorBuilder
	Remark(t task.Task) VisitorBuilder
	Script(t task.Task) VisitorBuilder
	Statement(t task.Task) VisitorBuilder
	WithContext(context.Context) Visitor
}

type visitorBuilder struct {
	call      task.Task
	print     task.Task
	remark    task.Task
	script    task.Task
	statement task.Task
}

func NewVisitor() VisitorBuilder {
	return &visitorBuilder{}
}

const (
	contextKey = "go-basic/Visitor"
)

func (b *visitorBuilder) WithContext(ctx context.Context) Visitor {
	v := &visitor{
		call:      b.call,
		print:     b.print,
		remark:    b.remark,
		script:    b.script,
		statement: b.statement,
	}
	v.ctx = context.WithValue(ctx, contextKey, v)
	return v
}

func FromContext(ctx context.Context) Visitor {
	return ctx.Value(contextKey).(Visitor)
}

func (b *visitorBuilder) Call(t task.Task) VisitorBuilder {
	b.call = t
	return b
}

func (b *visitorBuilder) Print(t task.Task) VisitorBuilder {
	b.print = t
	return b
}

func (b *visitorBuilder) Remark(t task.Task) VisitorBuilder {
	b.remark = t
	return b
}

func (b *visitorBuilder) Script(t task.Task) VisitorBuilder {
	b.script = t
	return b
}

func (b *visitorBuilder) Statement(t task.Task) VisitorBuilder {
	b.statement = t
	return b
}
