package script

import (
	"context"
	"github.com/alecthomas/participle/v2/lexer"
)

type VarDec struct {
	Pos lexer.Position

	ArrayDec  *ArrayDec  `parser:"  @@"`
	ScalarDec *ScalarDec `parser:"| @@"`
}

type ScalarDec struct {
	Pos lexer.Position

	Type string `parser:"@Type"`
	Name string `parser:"@Ident"`
}

func (s *ScalarDec) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, scalarDecKey, s)
}

func ScalarDecFromContext(ctx context.Context) *ScalarDec {
	return ctx.Value(scalarDecKey).(*ScalarDec)
}

type ArrayDec struct {
	Pos  lexer.Position
	Type string `parser:"@Type"`
	Name string `parser:"@Ident"`
	Size int    `parser:"\"[\" @Int \"]\""`
}

func (s *ArrayDec) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, arrayDecKey, s)
}

func ArrayDecFromContext(ctx context.Context) *ArrayDec {
	return ctx.Value(arrayDecKey).(*ArrayDec)
}
