package script

import (
	"context"
	"github.com/alecthomas/participle/v2/lexer"
)

type Script struct {
	Pos lexer.Position

	Include  []*Include `parser:"( @@"`
	FunDec   []*FuncDec `parser:"| @@)+"`
	Includes map[string]interface{}
}

func (s *Script) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, scriptKey, s)
}

func ScriptFromContext(ctx context.Context) *Script {
	return ctx.Value(scriptKey).(*Script)
}

type Include struct {
	Pos lexer.Position

	//Path *Path  `parser:"'(' @@ ')'"`
	Path []string `parser:"'include' ( @String (',' @String)* )"`
	//Path string `parser:"'!' 'include' @String "`
}

type Path struct {
	Pos lexer.Position

	Path string `parser:"@String"`
}
