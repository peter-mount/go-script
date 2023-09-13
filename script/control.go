package script

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type For struct {
	Pos lexer.Position

	Init      *Expression `parser:"'for' (@@)? ';'"`
	Condition *Expression `parser:"(@@)? ';'"`
	Increment *Expression `parser:"(@@)?"`
	Body      *Statement  `parser:"@@"`
}

// ForRange emulates go's "for i,v:=range expr {...}"
type ForRange struct {
	Pos lexer.Position

	Key        string      `parser:"'for' @Ident ','"` // index in range, _ to ignore
	Value      string      `parser:"@Ident"`           // value in range, _ to ignore
	Declare    bool        `parser:"@(':')?"`          // := to declare in local scope
	Expression *Expression `parser:" '=' 'range' @@"`
	Body       *Statement  `parser:"@@"`
}

type DoWhile struct {
	Pos lexer.Position

	Body      *Statement  `parser:"'do' @@"`
	Condition *Expression `parser:"'while' @@"`
}

type Repeat struct {
	Pos lexer.Position

	Body      *Statement  `parser:"'repeat' @@"`
	Condition *Expression `parser:"'until' @@"`
}

type While struct {
	Pos lexer.Position

	Condition *Expression `parser:"'while' @@"`
	Body      *Statement  `parser:"@@"`
}

type If struct {
	Pos lexer.Position

	Condition *Expression `parser:"'if' @@"`
	Body      *Statement  `parser:"@@"`
	Else      *Statement  `parser:"('else' @@)?"`
}

type Switch struct {
	Pos lexer.Position

	Expression *Expression   `parser:"'switch' (@@)? '{'"`
	Case       []*SwitchCase `parser:"(@@)+ "`
	Default    *Statement    `parser:"('default' ':' @@ )? '}'"`
}

type SwitchCase struct {
	Pos lexer.Position

	Expression []*SwitchCaseExpression `parser:"'case' (@@ (',' @@)*) ':'"`
	Statement  *Statement              `parser:"@@"`
}

type SwitchCaseExpression struct {
	Pos lexer.Position

	String     *string     `parser:"(@String"`
	Expression *Expression `parser:"| @@)"`
}
