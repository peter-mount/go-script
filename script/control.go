package script

import (
	"context"
	"github.com/alecthomas/participle/v2/lexer"
)

type ForStmt struct {
	Pos lexer.Position

	Init      *Expression `parser:"'for' (@@)? ';'"`
	Condition *Expression `parser:"(@@)? ';'"`
	Increment *Expression `parser:"(@@)?"`
	Body      *Statement  `parser:"@@"`
}

func (s *ForStmt) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, forKey, s)
}

func ForFromContext(ctx context.Context) *ForStmt {
	v := ctx.Value(forKey)
	if v != nil {
		return v.(*ForStmt)
	}
	return nil
}

// ForRange emulates go's "for i,v:=range expr {...}"
type ForRange struct {
	Pos lexer.Position

	Key        string      `parser:"'for' @Ident ','"` // index in range, _ to ignore
	Value      string      `parser:"@Ident"`           // value in range, _ to ignore
	Declare    bool        `parser:"@(':')?"`          // := to declare in local scope
	Expression *Expression `parser:" '=' 'range' @@"`
	Body       *Statement  `parser:"@@"`
}

func (s *ForRange) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, forRangeKey, s)
}

func ForRangeFromContext(ctx context.Context) *ForRange {
	v := ctx.Value(forRangeKey)
	if v != nil {
		return v.(*ForRange)
	}
	return nil
}

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
