package script

import (
	"context"
	"github.com/alecthomas/participle/v2/lexer"
)

type Call struct {
	Pos lexer.Position

	Name string        `@Ident`
	Args []*Expression `parser:"\"(\" ( @@ ( \",\" @@ )* )? \")\""`
}

func (s *Call) Accept(v Visitor) error { return nil } //v.VisitPrint(s) }

func (s *Call) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, callkey, s)
}

func CallFromContext(ctx context.Context) *Call {
	return ctx.Value(callkey).(*Call)
}
