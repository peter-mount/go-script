package script

import (
	"context"
	"github.com/alecthomas/participle/v2/lexer"
)

type Statement struct {
	Pos    lexer.Position
	Index  int     `parser:"@Number"`
	Remark *Remark `parser:"(   @@"`
	Print  *Print  `parser:"  | @@"`
	Call   *Call   `parser:"  | @@ ) EOL"`
}

func (s *Statement) Accept(v Visitor) error { return v.VisitStatement(s) }

func (s *Statement) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, statementKey, s)
}

func StatementFromContext(ctx context.Context) *Statement {
	return ctx.Value(statementKey).(*Statement)
}
