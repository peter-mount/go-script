package visitor

import (
	"context"
	"github.com/peter-mount/go-kernel/v2/util/task"
)

const (
	contextKey = "go-script/visitor"
)

type Builder interface {
	Addition(t task.Task) Builder
	// ArrayDec adds a task to invoke for each ArrayDec object
	ArrayDec(t task.Task) Builder
	// ArrayParamDec adds a task to invoke for each ArrayParamDec object
	ArrayParamDec(t task.Task) Builder
	Assignment(t task.Task) Builder
	// CallFunc adds a task to invoke for each CallFunc
	CallFunc(t task.Task) Builder
	Comparison(t task.Task) Builder
	Equality(t task.Task) Builder
	Expression(t task.Task) Builder
	ExpressionNoNest() Builder
	// FuncBody adds a task to invoke for each FuncBody object
	FuncBody(t task.Task) Builder
	// FuncDec adds a task to invoke for each FuncDec object
	FuncDec(t task.Task) Builder
	Multiplication(t task.Task) Builder
	Primary(t task.Task) Builder
	// ScalarDec adds a task to invoke for each ScalarDec object
	ScalarDec(t task.Task) Builder
	// ScalarParamDec adds a task to invoke for each ScalarParamDec object
	ScalarParamDec(t task.Task) Builder
	// Script adds a task to invoke for each Script object
	Script(t task.Task) Builder
	// Statement adds a task to invoke for each Statement object
	Statement(t task.Task) Builder
	// Statements adds a task to invoke for each Statements object
	Statements(t task.Task) Builder
	// StatementsNoNest tells the Visitor not to visit the Statement's within
	// a Statements object when it visits one.
	// This is useful when you want to process a Statements but want to handle the
	// content separately - e.g. Executor uses this.
	StatementsNoNest() Builder
	Unary(t task.Task) Builder
	// WithContext creates a Visitor bound to a Context
	WithContext(context.Context) Visitor
}

type builder struct {
	visitorCommon
}

func New() Builder {
	return &builder{}
}

func (b *builder) Addition(t task.Task) Builder {
	b.addition = b.addition.Then(t)
	return b
}

func (b *builder) ArrayDec(t task.Task) Builder {
	b.arrayDec = b.arrayDec.Then(t)
	return b
}

func (b *builder) ArrayParamDec(t task.Task) Builder {
	b.arrayParamDec = b.arrayParamDec.Then(t)
	return b
}

func (b *builder) Assignment(t task.Task) Builder {
	b.assignment = b.assignment.Then(t)
	return b
}

func (b *builder) CallFunc(t task.Task) Builder {
	b.callFunc = b.callFunc.Then(t)
	return b
}

func (b *builder) Comparison(t task.Task) Builder {
	b.comparison = b.comparison.Then(t)
	return b
}

func (b *builder) Equality(t task.Task) Builder {
	b.equality = b.equality.Then(t)
	return b
}

func (b *builder) Expression(t task.Task) Builder {
	b.expression = b.expression.Then(t)
	return b
}

func (b *builder) ExpressionNoNest() Builder {
	b.expressionNoNest = true
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

func (b *builder) Multiplication(t task.Task) Builder {
	b.multiplication = b.multiplication.Then(t)
	return b
}

func (b *builder) Primary(t task.Task) Builder {
	b.primary = b.primary.Then(t)
	return b
}

func (b *builder) ScalarDec(t task.Task) Builder {
	b.scalarDec = b.scalarDec.Then(t)
	return b
}

func (b *builder) ScalarParamDec(t task.Task) Builder {
	b.scalarParamDec = b.scalarParamDec.Then(t)
	return b
}

func (b *builder) Script(t task.Task) Builder {
	b.script = b.script.Then(t)
	return b
}

func (b *builder) Statements(t task.Task) Builder {
	b.statements = b.statements.Then(t)
	return b
}

func (b *builder) StatementsNoNest() Builder {
	b.statementsNoNest = true
	return b
}

func (b *builder) Statement(t task.Task) Builder {
	b.statement = b.statement.Then(t)
	return b
}

func (b *builder) Unary(t task.Task) Builder {
	b.unary = b.unary.Then(t)
	return b
}

func (b *builder) WithContext(ctx context.Context) Visitor {
	v := &visitor{visitorCommon: b.visitorCommon}
	v.ctx = context.WithValue(ctx, contextKey, v)
	return v
}