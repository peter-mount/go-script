package visitor

import (
	"context"
	"github.com/peter-mount/go-kernel/v2/util/task"
)

const (
	contextKey = "go-script/visitor"
)

type Builder interface {
	WithContext(context.Context) Visitor
	Script(t task.Task) Builder
	FuncDec(t task.Task) Builder
	FuncBody(t task.Task) Builder
	ArrayParamDec(t task.Task) Builder
	ScalarParamDec(t task.Task) Builder
	ArrayDec(t task.Task) Builder
	ScalarDec(t task.Task) Builder
}

type builder struct {
	script         task.Task
	funcDec        task.Task
	funcBody       task.Task
	arrayParamDec  task.Task
	scalarParamDec task.Task
	arrayDec       task.Task
	scalarDec      task.Task
}

func New() Builder {
	return &builder{}
}

func (b *builder) WithContext(ctx context.Context) Visitor {
	v := &visitor{
		script:         b.script,
		funcDec:        b.funcDec,
		funcBody:       b.funcBody,
		arrayParamDec:  b.arrayParamDec,
		scalarParamDec: b.scalarParamDec,
		arrayDec:       b.arrayDec,
		scalarDec:      b.scalarDec,
	}
	v.ctx = context.WithValue(ctx, contextKey, v)
	return v
}

func (b *builder) Script(t task.Task) Builder {
	b.script = b.script.Then(t)
	return b
}

func (b *builder) FuncDec(t task.Task) Builder {
	b.funcDec = b.funcDec.Then(t)
	return b
}

func (b *builder) FuncBody(t task.Task) Builder {
	b.funcBody = b.funcBody.Then(t)
	return b
}

func (b *builder) ArrayParamDec(t task.Task) Builder {
	b.arrayParamDec = b.arrayParamDec.Then(t)
	return b
}

func (b *builder) ScalarParamDec(t task.Task) Builder {
	b.scalarParamDec = b.scalarParamDec.Then(t)
	return b
}

func (b *builder) ArrayDec(t task.Task) Builder {
	b.arrayDec = b.arrayDec.Then(t)
	return b
}

func (b *builder) ScalarDec(t task.Task) Builder {
	b.scalarDec = b.scalarDec.Then(t)
	return b
}
