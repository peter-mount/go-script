package script

import (
	"context"
	"github.com/alecthomas/participle/v2/lexer"
)

type Print struct {
	Pos        lexer.Position
	Expression *Expression `parser:"\"print\" \"(\" @@ \")\""`
}

func (s *Print) Accept(v Visitor) error { return v.VisitPrint(s) }

func (s *Print) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, printKey, s)
}

func PrintFromContext(ctx context.Context) *Print {
	return ctx.Value(printKey).(*Print)
}
