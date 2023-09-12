package script

import (
	"context"
	"github.com/alecthomas/participle/v2/lexer"
)

type Expression struct {
	Pos lexer.Position

	Right *Assignment `parser:"@@"`
}

func (op *Expression) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, expressionKey, op)
}

func ExpressionFromContext(ctx context.Context) *Expression {
	return ctx.Value(expressionKey).(*Expression)
}

type Assignment struct {
	Pos lexer.Position

	Left    *Ternary `parser:"@@"`        // Expression or ident/reference to value to set
	Declare bool     `parser:"( @(':')?"` // := to declare in local scope, unset to use outer if already defined
	Op      string   `parser:"  @'='"`    // assign value
	Right   *Ternary `parser:"  @@ )?"`   // Expression to define value
}

func (op *Assignment) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, assignmentKey, op)
}

func AssignmentFromContext(ctx context.Context) *Assignment {
	return ctx.Value(assignmentKey).(*Assignment)
}

type Ternary struct {
	Pos lexer.Position

	Left  *Logic `parser:"@@"`
	True  *Logic `parser:"( '?' @@"`
	False *Logic `parser:"  ':' @@ )?"`
}

func (op *Ternary) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ternaryKey, op)
}

func TernaryFromContext(ctx context.Context) *Ternary {
	return ctx.Value(ternaryKey).(*Ternary)
}

type Logic struct {
	Pos lexer.Position

	Left  *Equality `parser:"@@"`
	Op    string    `parser:"[ @( '&' '&' | '|' '|' )"`
	Right *Logic    `parser:"  @@ ]"`
}

func (op *Logic) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, logicKey, op)
}

func LogicFromContext(ctx context.Context) *Logic {
	return ctx.Value(logicKey).(*Logic)
}

type Equality struct {
	Pos lexer.Position

	Left  *Comparison `parser:"@@"`
	Op    string      `parser:"[ @( '!' '=' | '=' '=' )"`
	Right *Equality   `parser:"  @@ ]"`
}

func (op *Equality) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, equalityKey, op)
}

func EqualityFromContext(ctx context.Context) *Equality {
	return ctx.Value(equalityKey).(*Equality)
}

type Comparison struct {
	Pos lexer.Position

	Left  *Addition   `parser:"@@"`
	Op    string      `parser:"[ @( '>' '=' | '>' | '<' '=' | '<' )"`
	Right *Comparison `parser:"  @@ ]"`
}

func (op *Comparison) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, comparisonKey, op)
}

func ComparisonFromContext(ctx context.Context) *Comparison {
	return ctx.Value(comparisonKey).(*Comparison)
}

type Addition struct {
	Pos lexer.Position

	Left  *Multiplication `parser:"@@"`
	Op    string          `parser:"[ @( '+' | '-' )"`
	Right *Addition       `parser:"  @@ ]"`
}

func (op *Addition) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, additionKey, op)
}

func AdditionFromContext(ctx context.Context) *Addition {
	return ctx.Value(additionKey).(*Addition)
}

type Multiplication struct {
	Pos lexer.Position

	Left  *Unary          `parser:"@@"`
	Op    string          `parser:"[ @( '*' | '/' | '%' )"`
	Right *Multiplication `parser:"  @@ ]"`
}

func (op *Multiplication) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, multiplicationKey, op)
}

func MultiplicationFromContext(ctx context.Context) *Multiplication {
	return ctx.Value(multiplicationKey).(*Multiplication)
}

type Unary struct {
	Pos lexer.Position

	Op    string   `parser:"  ( @( '!' | '-' )"`
	Left  *Unary   `parser:"    @@ )"`
	Right *Primary `parser:"| @@"`
}

func (op *Unary) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, unaryKey, op)
}

func UnaryFromContext(ctx context.Context) *Unary {
	return ctx.Value(unaryKey).(*Unary)
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
	Ident         *Ident      `parser:"  |   @@ "`
	PointOp       string      `parser:"    ) [ @Period"`
	Pointer       *Primary    `parser:"      @@] )"`
}

func (op *Primary) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, primaryKey, op)
}

func PrimaryFromContext(ctx context.Context) *Primary {
	return ctx.Value(primaryKey).(*Primary)
}

type Ident struct {
	Pos lexer.Position

	Ident string        `parser:"@Ident"`
	Index []*Expression `parser:"[ ('[' @@ ']')+ ]"`
}

// KeyValue is "string": expression
type KeyValue struct {
	Pos lexer.Position

	Key   string      `parser:"@String"`
	Value *Expression `parser:"':' @@"`
}
