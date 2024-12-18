package script

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type Script struct {
	Pos lexer.Position

	Import   []*Import  `parser:"( @@"`
	Include  []*Include `parser:"| @@"`
	FunDec   []*FuncDec `parser:"| @@)+"`
	Includes map[string]interface{}
}

type Import struct {
	Pos lexer.Position

	Packages []*ImportPackage `parser:"'import' '(' ( @@+ ) ')'"`
}

type ImportPackage struct {
	Pos lexer.Position

	As   string `parser:"( @Ident )?"`
	Name string `parser:"@String"`
}

type Include struct {
	Pos lexer.Position

	Path []string `parser:"'include' ( @String (',' @String)* )"`
}
