package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/peter-mount/go-script/script"
)

var (
	scriptLexer = lexer.MustSimple([]lexer.SimpleRule{
		{"hashComment", `#.*`},
		{"comment", `//.*|/\*.*?\*/`},
		{"whitespace", `\s+`},
		{"Ident", `\b(([a-zA-Z_][a-zA-Z0-9_]*)(\.([a-zA-Z_][a-zA-Z0-9_]*))*)\b`},
		{"Punct", `[-,()*/+%{};&!=:<>]|\[|\]`},
		{"Number", `[-+]?((\d*)?\.\d+|\d+\.(\d*)?)`},
		{"Int", `[-+]?\d+`},
		{"String", `"(\\"|[^"])*"`},
	})

	scriptParser = participle.MustBuild[script.Script](
		participle.Lexer(scriptLexer),
		participle.UseLookahead(2),
		participle.Unquote("String"),
	)
)
