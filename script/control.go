package script

import (
	"context"
	"github.com/alecthomas/participle/v2/lexer"
)

type For struct {
	Pos lexer.Position

	Init      *Expression `parser:"'for' (@@)? ';'"`
	Condition *Expression `parser:"(@@)? ';'"`
	Increment *Expression `parser:"(@@)?"`
	Body      *Statement  `parser:"@@"`
}

func (s *For) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, forKey, s)
}

func ForFromContext(ctx context.Context) *For {
	v := ctx.Value(forKey)
	if v != nil {
		return v.(*For)
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

type While struct {
	Pos lexer.Position

	Condition *Expression `parser:"'while' @@"`
	Body      *Statement  `parser:"@@"`
}

func (s *While) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, whileKey, s)
}

func WhileFromContext(ctx context.Context) *While {
	v := ctx.Value(whileKey)
	if v != nil {
		return v.(*While)
	}
	return nil
}

type If struct {
	Pos lexer.Position

	Condition *Expression `parser:"'if' @@"`
	Body      *Statement  `parser:"@@"`
	Else      *Statement  `parser:"('else' @@)?"`
}

func (s *If) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ifKey, s)
}

func IfFromContext(ctx context.Context) *If {
	v := ctx.Value(ifKey)
	if v != nil {
		return v.(*If)
	}
	return nil
}
