package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/peter-mount/go-script/script"
)

var (
	scriptLexer = lexer.MustSimple([]lexer.SimpleRule{
		{"comment", `//.*|/\*.*?\*/`},
		{"whitespace", `\s+`},

		{"Type", `\b(int|char)\b`},
		{"Ident", `\b([a-zA-Z_][a-zA-Z0-9_]*)\b`},
		{"Punct", `[-,()*/+%{};&!=:<>]|\[|\]`},
		{"Number", `[-+]?((\d*)?\.\d+|\d+\.(\d*)?)`},
		{"Int", `[-+]?\d+`},
		{"String", `"(\\"|[^"])*"`},
		//{"EOS", `[;\n\r]+`},
		//{"EOL", `[\n\r]+`},
	})

	scriptParser = participle.MustBuild[script.Script](
		participle.Lexer(scriptLexer),
		participle.UseLookahead(2),
		//participle.CaseInsensitive("Ident"),
		participle.Unquote("String"),
		//participle.Elide("whitespace"),
	)
)
