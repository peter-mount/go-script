package script

import (
	"context"
	"github.com/alecthomas/participle/v2/lexer"
)

type Script struct {
	Pos lexer.Position

	TopDec []*TopDec `parser:"@@*"`
}

func (s *Script) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, scriptKey, s)
}

func ScriptFromContext(ctx context.Context) *Script {
	return ctx.Value(scriptKey).(*Script)
}

type TopDec struct {
	Pos lexer.Position

	FunDec *FuncDec `parser:"@@"`
}
