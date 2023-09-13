package script

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type Script struct {
	Pos lexer.Position

	Include  []*Include `parser:"( @@"`
	FunDec   []*FuncDec `parser:"| @@)+"`
	Includes map[string]interface{}
}

type Include struct {
	Pos lexer.Position

	Path []string `parser:"'include' ( @String (',' @String)* )"`
}
