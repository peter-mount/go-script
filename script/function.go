package script

import (
	"context"
	"github.com/alecthomas/participle/v2/lexer"
)

type FuncDec struct {
	Pos lexer.Position

	Name       string      `parser:"@Ident"`
	Parameters []string    `parser:"'(' (@Ident (',' @Ident)*)? ')'"`
	FunBody    *Statements `parser:"@@"`
}

func (s *FuncDec) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, funcDecKey, s)
}

func FuncDecFromContext(ctx context.Context) *FuncDec {
	return ctx.Value(funcDecKey).(*FuncDec)
}

type Return struct {
	Pos lexer.Position

	Result *Expression `parser:"'return' @@?"`
}

func (s *Return) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, returnKey, s)
}

func ReturnFromContext(ctx context.Context) *Return {
	return ctx.Value(returnKey).(*Return)
}

type CallFunc struct {
	Pos lexer.Position

	Name       string         `parser:"@Ident"`
	Parameters *ParameterList `parser:"'(' @@? ')'"`
}

type ParameterList struct {
	Args     []*Expression `parser:"(@@ (',' @@)*) "`
	Variadic bool          `parser:"    @('...')?"`
}

func (s *CallFunc) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, callFuncKey, s)
}

func CallFuncFromContext(ctx context.Context) *CallFunc {
	return ctx.Value(callFuncKey).(*CallFunc)
}
