package script

import (
	"context"
	"github.com/alecthomas/participle/v2/lexer"
)

type Script struct {
	Pos        lexer.Position
	Statements []*Statement `parser:"@@*"`
	Table      map[int]*Statement
}

func (s *Script) Accept(v Visitor) error { return v.VisitScript(s) }

func (s *Script) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, scriptKey, s)
}

func ScriptFromContext(ctx context.Context) *Script {
	return ctx.Value(scriptKey).(*Script)
}
