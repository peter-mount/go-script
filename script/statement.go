package script

import (
	"context"
	"github.com/alecthomas/participle/v2/lexer"
)

type Statements struct {
	Pos    lexer.Position
	Parent *Statements // Parent when nested

	Statements []*Statement `parser:"'{' @@* '}'"`
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
	Next   *Statement  // Next statement within Statements

	Break    *Break    `parser:"  @@"`
	IfStmt   *If       `parser:"| @@"`
	For      *For      `parser:"| @@"`
	ForRange *ForRange `parser:"| @@"`
	Repeat   *Repeat   `parser:"| @@"`
	Return   *Return   `parser:"| @@"`
	While    *While    `parser:"| @@"`

	// Try is after the main block as it's a bit more complex,
	// so it's better to place it here after the statements
	// when in the railroad diagrams.
	//
	// Otherwise, this could be in the above block
	Try *Try `parser:"| @@"`

	// These must be at the end
	Block      *Statements `parser:"| @@"`
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
