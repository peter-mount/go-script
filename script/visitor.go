package script

import (
	"context"
	"github.com/peter-mount/go-kernel/v2/util/task"
)

type Visitor interface {
	VisitBlock(block *Block) error
	VisitCall(call *Call) error
	VisitPrint(print *Print) error
	VisitScript(script *Script) error
	VisitStatement(statement *Statement) error
}

type Visitable interface {
	Accept(v Visitor) error
}

type visitor struct {
	ctx       context.Context
	block     task.Task
	call      task.Task
	print     task.Task
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

func (v *visitor) VisitBlock(block *Block) error {
	return v.visit(block.WithContext, func() error {
		if err := v.script.Do(v.ctx); err != nil {
			return err
		}

		// FIXME use PC not sequential!
		for _, statement := range block.Statements {
			if err := statement.Accept(v); err != nil {
				return err
			}
		}
		return nil
	})
}

func (v *visitor) VisitCall(call *Call) error {
	return v.visitTask(call.WithContext, v.call)
}

func (v *visitor) VisitPrint(print *Print) error {
	return v.visitTask(print.WithContext, v.print)
}

func (v *visitor) VisitScript(s *Script) error {
	return v.visit(s.WithContext, func() error {
		if err := v.script.Do(v.ctx); err != nil {
			return err
		}

		// FIXME use PC not sequential!
		/*for _, block := range s.Blocks {
			if err := block.Accept(v); err != nil {
				return err
			}
		}*/
		return nil
	})
}

func (v *visitor) VisitStatement(s *Statement) error {
	return v.visit(s.WithContext, func() error {
		if err := v.statement.Do(v.ctx); err != nil {
			return err
		}

		switch {
		/*		case s.Call != nil:
					return s.Call.Accept(v)
				case s.Print != nil:
					return s.Print.Accept(v)
		*/
		}

		// TODO raise error here?
		return nil
	})
}

type VisitorBuilder interface {
	Block(t task.Task) VisitorBuilder
	Call(t task.Task) VisitorBuilder
	Print(t task.Task) VisitorBuilder
	Script(t task.Task) VisitorBuilder
	Statement(t task.Task) VisitorBuilder
	WithContext(context.Context) Visitor
}

type visitorBuilder struct {
	block     task.Task
	call      task.Task
	print     task.Task
	script    task.Task
	statement task.Task
}

func NewVisitor() VisitorBuilder {
	return &visitorBuilder{}
}

const (
	contextKey = "go-script/Visitor"
)

func (b *visitorBuilder) WithContext(ctx context.Context) Visitor {
	v := &visitor{
		block:     b.block,
		call:      b.call,
		print:     b.print,
		script:    b.script,
		statement: b.statement,
	}
	v.ctx = context.WithValue(ctx, contextKey, v)
	return v
}

func FromContext(ctx context.Context) Visitor {
	return ctx.Value(contextKey).(Visitor)
}

func (b *visitorBuilder) Block(t task.Task) VisitorBuilder {
	b.block = b.block.Then(t)
	return b
}

func (b *visitorBuilder) Call(t task.Task) VisitorBuilder {
	b.call = b.call.Then(t)
	return b
}

func (b *visitorBuilder) Print(t task.Task) VisitorBuilder {
	b.print = b.print.Then(t)
	return b
}

func (b *visitorBuilder) Script(t task.Task) VisitorBuilder {
	b.script = b.script.Then(t)
	return b
}

func (b *visitorBuilder) Statement(t task.Task) VisitorBuilder {
	b.statement = b.statement.Then(t)
	return b
}
