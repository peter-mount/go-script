package script

import (
	"context"
	"github.com/alecthomas/participle/v2/lexer"
)

type FuncDec struct {
	Pos lexer.Position

	ReturnType string       `parser:"@(Type | 'void')?"`
	Name       string       `parser:"@Ident"`
	Parameters []*Parameter `parser:"'(' ((@@ (',' @@)*) | 'void' )? ')'"`
	FunBody    *FuncBody    `parser:"(';' | '{' @@ '}')"`
}

func (s *FuncDec) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, funcDecKey, s)
}

func FuncDecFromContext(ctx context.Context) *FuncDec {
	return ctx.Value(funcDecKey).(*FuncDec)
}

type FuncBody struct {
	Pos lexer.Position

	Locals     []*VarDec   `parser:"(@@ ';')*"`
	Statements *Statements `parser:"@@"`
}

func (s *FuncBody) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, funcBodyKey, s)
}

func FuncBodyFromContext(ctx context.Context) *FuncBody {
	return ctx.Value(funcBodyKey).(*FuncBody)
}

type Parameter struct {
	Pos lexer.Position

	Array  *ArrayParameter `parser:"  @@"`
	Scalar *ScalarDec      `parser:"| @@"`
}

type ArrayParameter struct {
	Pos lexer.Position

	Type string `parser:"@Type"`
	Name string `parser:"@Ident '[' ']'"`
}

func (s *ArrayParameter) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, arrayParameterKey, s)
}

func ArrayParameterFromContext(ctx context.Context) *ArrayParameter {
	return ctx.Value(arrayParameterKey).(*ArrayParameter)
}

type ReturnStmt struct {
	Pos lexer.Position

	Result *Expression `parser:"'return' @@?"`
}

type CallFunc struct {
	Pos lexer.Position

	Name string        `parser:"@Ident"`
	Args []*Expression `parser:"'(' (@@ (',' @@)*)? ')'"`
}

func (s *CallFunc) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, callFuncKey, s)
}

func CallFuncFromContext(ctx context.Context) *CallFunc {
	return ctx.Value(callFuncKey).(*CallFunc)
}
