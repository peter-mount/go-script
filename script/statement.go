package script

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type Statements struct {
	Pos lexer.Position

	Statements []*Statement `parser:"'{' @@* '}'"`
}

type Statement struct {
	Pos  lexer.Position
	Next *Statement // Next statement within Statements block

	Break    bool      `parser:"  @'break'"`
	Continue bool      `parser:"| @'continue'"`
	DoWhile  *DoWhile  `parser:"| @@"`
	IfStmt   *If       `parser:"| @@"`
	For      *For      `parser:"| @@"`
	ForRange *ForRange `parser:"| @@"`
	Repeat   *Repeat   `parser:"| @@"`
	Return   *Return   `parser:"| @@"`
	Switch   *Switch   `parser:"| @@"`
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
