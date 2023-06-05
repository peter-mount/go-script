package script

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type Expression struct {
	Pos lexer.Position

	Assignment *Assignment `parser:"@@"`
}

type Assignment struct {
	Pos lexer.Position

	Equality *Equality `parser:"@@"`
	Op       string    `parser:"( @'='"`
	Next     *Equality `parser:"  @@ )?"`
}

type Equality struct {
	Pos lexer.Position

	Comparison *Comparison `parser:"@@"`
	Op         string      `parser:"[ @( '!' '=' | '=' '=' )"`
	Next       *Equality   `parser:"  @@ ]"`
}

type Comparison struct {
	Pos lexer.Position

	Addition *Addition   `parser:"@@"`
	Op       string      `parser:"[ @( '>' '=' | '>' | '<' '=' | '<' )"`
	Next     *Comparison `parser:"  @@ ]"`
}

type Addition struct {
	Pos lexer.Position

	Multiplication *Multiplication `parser:"@@"`
	Op             string          `parser:"[ @( '-' | '+' )"`
	Next           *Addition       `parser:"  @@ ]"`
}

type Multiplication struct {
	Pos lexer.Position

	Unary *Unary          `parser:"@@"`
	Op    string          `parser:"[ @( '/' | '*' )"`
	Next  *Multiplication `parser:"  @@ ]"`
}

type Unary struct {
	Pos lexer.Position

	Op      string   `parser:"  ( @( '!' | '-' )"`
	Unary   *Unary   `parser:"    @@ )"`
	Primary *Primary `parser:"| @@"`
}

type Primary struct {
	Pos lexer.Position

	Float         *float64    `parser:"  @Number"`
	Integer       *int        `parser:"| @Int"`
	String        *string     `parser:"| @String"`
	ArrayIndex    *ArrayIndex `parser:"| @@"`
	CallFunc      *CallFunc   `parser:"| @@"`
	Ident         string      `parser:"| @Ident"`
	SubExpression *Expression `parser:"| '(' @@ ')' "`
}

type ArrayIndex struct {
	Pos lexer.Position

	Ident string        `parser:"@Ident"`
	Index []*Expression `parser:"('[' @@ ']')+"`
}
