package script

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type Expression struct {
	Pos lexer.Position

	Right *Assignment `parser:"@@"`
}

type Assignment struct {
	Pos lexer.Position

	Left        *Ternary    `parser:"@@"`                        // Expression or ident/reference to value to set
	AugmentedOp *string     `parser:"( @('+'|'-'|'*'|'/'|'%')?"` // Operation to perform on the result
	Declare     bool        `parser:"  @(':')?"`                 // := to declare in local scope, unset to use outer if already defined
	Op          string      `parser:"  @'='"`                    // assign value
	Right       *Assignment `parser:"  @@ )?"`                   // Expression to define value
}

type Ternary struct {
	Pos lexer.Position

	Left  *Logic `parser:"@@"`
	True  *Logic `parser:"( '?' @@"`
	False *Logic `parser:"  ':' @@ )?"`
}

type Logic struct {
	Pos lexer.Position

	Left  *Equality `parser:"@@"`
	Op    string    `parser:"[ @( '&' '&' | '|' '|' )"`
	Right *Logic    `parser:"  @@ ]"`
}

type Equality struct {
	Pos lexer.Position

	Left  *Comparison `parser:"@@"`
	Op    string      `parser:"[ @( '!' '=' | '=' '=' )"`
	Right *Equality   `parser:"  @@ ]"`
}

type Comparison struct {
	Pos lexer.Position

	Left  *Addition   `parser:"@@"`
	Op    string      `parser:"[ @( '>' '=' | '>' | '<' '=' | '<' )"`
	Right *Comparison `parser:"  @@ ]"`
}

type Addition struct {
	Pos lexer.Position

	Left  *Multiplication `parser:"@@"`
	Op    string          `parser:"[ @( '+' | '-' )"`
	Right *Addition       `parser:"  @@ ]"`
}

type Multiplication struct {
	Pos lexer.Position

	Left  *Unary          `parser:"@@"`
	Op    string          `parser:"[ @( '*' | '/' | '%' )"`
	Right *Multiplication `parser:"  @@ ]"`
}

type Unary struct {
	Pos lexer.Position

	Op    string   `parser:"  ( @( '!' | '-' )"`
	Left  *Primary `parser:"    @@ )"`
	Right *Primary `parser:"| @@"`
}

type Primary struct {
	Pos lexer.Position

	Float         *float64    `parser:"( @Number"`
	Integer       *int        `parser:"  | @Int"`
	KeyValue      *KeyValue   `parser:"  | @@"`
	String        *string     `parser:"  | @String"`
	Null          bool        `parser:"  | @'null'"`
	Nil           bool        `parser:"  | @'nil'"`
	True          bool        `parser:"  | @'true'"`
	False         bool        `parser:"  | @'false'"`
	SubExpression *Expression `parser:"  | '(' @@ ')' "`
	CallFunc      *CallFunc   `parser:"  | ( @@"`
	Ident         *Ident      `parser:"    | @@ "`
	PointOp       string      `parser:"    ) [ @Period"`
	Pointer       *Primary    `parser:"      @@] )"`
}

type Ident struct {
	Pos lexer.Position

	PreIncDec  *IncDec       `parser:"(@@?)"`
	Ident      string        `parser:"@Ident"`
	PostIncDec *IncDec       `parser:"(@@?)"`
	Index      []*Expression `parser:"[ ('[' @@ ']')+ ]"`
}

type IncDec struct {
	Pos lexer.Position

	Decrement bool `parser:"( @('-' '-')"`
	Increment bool `parser:"  | @('+' '+') )"`
}

// IsPreIncDec returns true if --ident or ++ident but no array indices
func (i *Ident) IsPreIncDec() bool {
	return i != nil && i.PreIncDec != nil && len(i.Index) == 0
}

// IsPostIncDec returns true if ident-- or ident++ but no array indices
func (i *Ident) IsPostIncDec() bool {
	return i != nil && i.PostIncDec != nil && len(i.Index) == 0
}

// IsIndexed returns true if ident[...] but no pre or post incDec
func (i *Ident) IsIndexed() bool {
	return i != nil && len(i.Index) > 0 && !(i.IsPreIncDec() || i.IsPostIncDec())
}

// KeyValue is "string": expression
type KeyValue struct {
	Pos lexer.Position

	Key   string      `parser:"@String"`
	Value *Expression `parser:"':' @@"`
}
