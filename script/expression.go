package script

import (
	"context"
	"fmt"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/peter-mount/go-script/calculator"
)

type Expression struct {
	Pos lexer.Position

	Assignment *Assignment `parser:"@@"`
}

func (op *Expression) Calculate(ctx context.Context) error {
	calc := calculator.FromContext(ctx)
	if calc != nil && op.Assignment != nil {
		// This ensures the Calculator has a fresh stack for each Expression
		v, err := calc.Exec(op.Assignment, ctx)
		if err != nil {
			return err
		}
		if v != nil {
			calc.Push(v)
		}
	}
	return nil
}

type Assignment struct {
	Pos lexer.Position

	Left  *Equality `parser:"@@"`
	Op    string    `parser:"( @'='"`
	Right *Equality `parser:"  @@ )?"`
}

func (op *Assignment) Calculate(ctx context.Context) error {
	if op.Op == "=" {
		return fmt.Errorf("assignment not implemented")
	}

	return calculator.DoCalculation(op.Left, ctx)
}

type Equality struct {
	Pos lexer.Position

	Left  *Comparison `parser:"@@"`
	Op    string      `parser:"[ @( '!' '=' | '=' '=' )"`
	Right *Equality   `parser:"  @@ ]"`
}

func (op *Equality) Calculate(ctx context.Context) error {
	return calculator.CalculateOp2(ctx, op.Op, op.Left, op.Right)
}

type Comparison struct {
	Pos lexer.Position

	Left  *Addition   `parser:"@@"`
	Op    string      `parser:"[ @( '>' '=' | '>' | '<' '=' | '<' )"`
	Right *Comparison `parser:"  @@ ]"`
}

func (op *Comparison) Calculate(ctx context.Context) error {
	return calculator.CalculateOp2(ctx, op.Op, op.Left, op.Right)
}

type Addition struct {
	Pos lexer.Position

	Left  *Multiplication `parser:"@@"`
	Op    string          `parser:"[ @( '-' | '+' )"`
	Right *Addition       `parser:"  @@ ]"`
}

func (op *Addition) Calculate(ctx context.Context) error {
	return calculator.CalculateOp2(ctx, op.Op, op.Left, op.Right)
}

type Multiplication struct {
	Pos lexer.Position

	Left  *Unary          `parser:"@@"`
	Op    string          `parser:"[ @( '/' | '*' )"`
	Right *Multiplication `parser:"  @@ ]"`
}

func (op *Multiplication) Calculate(ctx context.Context) error {
	return calculator.CalculateOp2(ctx, op.Op, op.Left, op.Right)
}

type Unary struct {
	Pos lexer.Position

	Op      string   `parser:"  ( @( '!' | '-' )"`
	Unary   *Unary   `parser:"    @@ )"`
	Primary *Primary `parser:"| @@"`
}

func (op *Unary) Calculate(ctx context.Context) error {
	// TODO how to handle? not a CalculateOp
	switch {
	case op.Op == "!":
		// Negate
		return fmt.Errorf("not '!' not implemented")
	case op.Op == "-":
		// negate
		return fmt.Errorf("negate '-' not implemented")
	case op.Primary != nil:
		return calculator.DoCalculation(op.Primary, ctx)
	default:
		return nil
	}
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

func (op *Primary) Calculate(ctx context.Context) error {
	calc := calculator.FromContext(ctx)
	if calc != nil {
		switch {
		case op.Float != nil:
			calc.Push(*op.Float)
		case op.Integer != nil:
			calc.Push(*op.Integer)
		case op.String != nil:
			calc.Push(*op.String)
		case op.Ident != "":
		//return *op.Float, nil
		case op.SubExpression != nil:
			return calculator.DoCalculation(op.SubExpression, ctx)
		}
	}
	return nil
}

type ArrayIndex struct {
	Pos lexer.Position

	Ident string        `parser:"@Ident"`
	Index []*Expression `parser:"('[' @@ ']')+"`
}
