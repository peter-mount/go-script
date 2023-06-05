package script

import (
	"context"
	"github.com/alecthomas/participle/v2/lexer"
)

type Statements struct {
	Pos    lexer.Position
	Parent *Statements // Parent when nested

	Statements []*Statement `parser:"@@*"`
}

func (s *Statements) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, statementsKey, s)
}

func StatementsFromContext(ctx context.Context) *Statements {
	v := ctx.Value(statementsKey)
	if v != nil {
		return v.(*Statements)
	}
	return nil
}

type Statement struct {
	Pos    lexer.Position
	Parent *Statements // Parent when nested

	IfStmt     *IfStmt     `parser:"  @@"`
	ReturnStmt *ReturnStmt `parser:"| @@"`
	WhileStmt  *WhileStmt  `parser:"| @@"`
	Block      *Statements `parser:"| '{' @@ '}'"`
	Expression *Expression `parser:"| @@"`
	Empty      bool        `parser:"| @';'"`
}

func (s *Statement) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, statementKey, s)
}

func StatementFromContext(ctx context.Context) *Statement {
	v := ctx.Value(statementKey)
	if v != nil {
		return v.(*Statement)
	}
	return nil
}
