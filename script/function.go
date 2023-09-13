package script

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type FuncDec struct {
	Pos lexer.Position

	Name       string      `parser:"@Ident"`
	Parameters []string    `parser:"'(' (@Ident (',' @Ident)*)? ')'"`
	FunBody    *Statements `parser:"@@"`
}

type Return struct {
	Pos lexer.Position

	Result *Expression `parser:"'return' @@?"`
}

type CallFunc struct {
	Pos lexer.Position

	Name       string         `parser:"@Ident"`
	Parameters *ParameterList `parser:"'(' @@? ')'"`
}

type ParameterList struct {
	Args     []*Expression `parser:"(@@ (',' @@)*) "`
	Variadic bool          `parser:"    @('...')?"`
}
