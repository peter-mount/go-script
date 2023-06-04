package script

import "github.com/alecthomas/participle/v2/lexer"

type WhileStmt struct {
	Pos lexer.Position

	Condition *Expression `"while" "(" @@ ")"`
	Body      *Statement  `@@`
}

type IfStmt struct {
	Pos lexer.Position

	Condition *Expression `"if" "(" @@ ")"`
	Body      *Statement  `@@`
	Else      *Statement  `("else" @@)?`
}
