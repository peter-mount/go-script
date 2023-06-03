package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/peter-mount/go-basic/script"
	"io"
	"os"
)

func Parse(fileName string, r io.Reader, opts ...participle.ParseOption) (*script.Script, error) {
	return parser.Parse(fileName, r, opts...)
}

func ParseBytes(fileName string, b []byte, opts ...participle.ParseOption) (*script.Script, error) {
	return parser.ParseBytes(fileName, b, opts...)
}

func ParseString(fileName, s string, opts ...participle.ParseOption) (*script.Script, error) {
	return parser.ParseString(fileName, s, opts...)
}

func ParseFile(fileName string, opts ...participle.ParseOption) (*script.Script, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parser.Parse(fileName, f, opts...)
}
