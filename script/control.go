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

type DoWhile struct {
	Pos lexer.Position

	Body      *Statement  `parser:"'do' @@"`
	Condition *Expression `parser:"'while' @@"`
}

func (s *DoWhile) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, doWhileKey, s)
}

func DoWhileFromContext(ctx context.Context) *DoWhile {
	v := ctx.Value(doWhileKey)
	if v != nil {
		return v.(*DoWhile)
	}
	return nil
}

type Repeat struct {
	Pos lexer.Position

	Body      *Statement  `parser:"'repeat' @@"`
	Condition *Expression `parser:"'until' @@"`
}

func (s *Repeat) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, repeatKey, s)
}

func RepeatFromContext(ctx context.Context) *Repeat {
	v := ctx.Value(repeatKey)
	if v != nil {
		return v.(*Repeat)
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

type Switch struct {
	Pos lexer.Position

	Expression *Expression   `parser:"'switch' (@@)? '{'"`
	Case       []*SwitchCase `parser:"(@@)+ "`
	Default    *Statement    `parser:"('default' ':' @@ )? '}'"`
}

type SwitchCase struct {
	Pos lexer.Position

	String     *string     `parser:"'case' ( @String"` // For some reason we need to check String specifically
	Expression *Expression `parser:"| @@) ':'"`        // otherwise the parser can't handle it in expression.
	Statement  *Statement  `parser:"@@"`               // But it works as expression elsewhere... odd
}

func (s *Switch) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, switchKey, s)
}

func SwitchFromContext(ctx context.Context) *Switch {
	v := ctx.Value(switchKey)
	if v != nil {
		return v.(*Switch)
	}
	return nil
}
