package script

import (
	"context"
	"github.com/alecthomas/participle/v2/lexer"
)

type WhileStmt struct {
	Pos lexer.Position

	Condition *Expression `parser:"'while' @@"`
	Body      *Statement  `parser:"@@"`
}

func (s *WhileStmt) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, whileKey, s)
}

func WhileFromContext(ctx context.Context) *WhileStmt {
	v := ctx.Value(whileKey)
	if v != nil {
		return v.(*WhileStmt)
	}
	return nil
}

type IfStmt struct {
	Pos lexer.Position

	Condition *Expression `parser:"'if' @@"`
	Body      *Statement  `parser:"@@"`
	Else      *Statement  `parser:"('else' @@)?"`
}

func (s *IfStmt) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ifKey, s)
}

func IfFromContext(ctx context.Context) *IfStmt {
	v := ctx.Value(ifKey)
	if v != nil {
		return v.(*IfStmt)
	}
	return nil
}
