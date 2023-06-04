package script

import "github.com/alecthomas/participle/v2/lexer"

type VarDec struct {
	Pos lexer.Position

	ArrayDec  *ArrayDec  `parser:"  @@"`
	ScalarDec *ScalarDec `parser:"| @@"`
}

type ScalarDec struct {
	Pos lexer.Position

	Type string `parser:"@Type"`
	Name string `parser:"@Ident"`
}

type ArrayDec struct {
	Pos  lexer.Position
	Type string `parser:"@Type"`
	Name string `parser:"@Ident"`
	Size int    `parser:"\"[\" @Int \"]\""`
}
