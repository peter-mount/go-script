package script

import (
	"context"
	"github.com/alecthomas/participle/v2/lexer"
)

type Statements struct {
	Pos        lexer.Position
	Statements []*Statement `parser:"@@*"`
}

type Statement struct {
	Pos lexer.Position

	IfStmt     *IfStmt     `parser:"  @@"`
	ReturnStmt *ReturnStmt `parser:"| @@"`
	WhileStmt  *WhileStmt  `parser:"| @@"`
	Block      *Statements `parser:"| \"{\" @@ \"}\""`
	Expression *Expression `parser:"| @@"`
	Empty      bool        `parser:"| @\";\""`
}

func (s *Statement) Accept(v Visitor) error { return v.VisitStatement(s) }

func (s *Statement) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, statementKey, s)
}

func StatementFromContext(ctx context.Context) *Statement {
	return ctx.Value(statementKey).(*Statement)
}
