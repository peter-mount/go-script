package script

import "github.com/alecthomas/participle/v2/lexer"

type FunDec struct {
	Pos lexer.Position

	ReturnType string       `parser:"@(Type | \"void\")?"`
	Name       string       `parser:"@Ident"`
	Parameters []*Parameter `parser:"\"(\" ((@@ (\",\" @@)*) | \"void\" )? \")\""`
	FunBody    *FunBody     `parser:"(\";\" | \"{\" @@ \"}\")"`
}

type FunBody struct {
	Pos lexer.Position

	Locals     []*VarDec   `parser:"(@@ \";\")*"`
	Statements *Statements `parser:"@@"`
}

type Parameter struct {
	Pos lexer.Position

	Array  *ArrayParameter `parser:"  @@"`
	Scalar *ScalarDec      `parser:"| @@"`
}

type ArrayParameter struct {
	Pos lexer.Position

	Type  string `parser:"@Type"`
	Ident string `parser:"@Ident \"[\" \"]\""`
}

type ReturnStmt struct {
	Pos lexer.Position

	Result *Expression `parser:"\"return\" @@?"`
}

type CallFunc struct {
	Pos lexer.Position

	Ident string        `parser:"@Ident"`
	Index []*Expression `parser:"\"(\" (@@ (\",\" @@)*)? \")\""`
}
