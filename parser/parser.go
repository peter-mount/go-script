package parser

import (
	"fmt"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/peter-mount/go-script/errors"
	"github.com/peter-mount/go-script/script"
	"io"
	"os"
	"path/filepath"
)

type Parser interface {
	Parse(fileName string, r io.Reader, opts ...participle.ParseOption) (*script.Script, error)
	ParseBytes(fileName string, b []byte, opts ...participle.ParseOption) (*script.Script, error)
	ParseString(fileName, src string, opts ...participle.ParseOption) (*script.Script, error)
	ParseFile(fileName string, opts ...participle.ParseOption) (*script.Script, error)
	IncludePath(s string) error
	EBNF() string
}

type defaultParser struct {
	lexer       *lexer.StatefulDefinition
	parser      *participle.Parser[script.Script]
	includePath []string
}

func New() Parser {
	return &defaultParser{
		lexer:  scriptLexer,
		parser: scriptParser,
	}
}

func (p *defaultParser) EBNF() string {
	return p.parser.String()
}

func (p *defaultParser) IncludePath(s string) error {
	fi, err := os.Stat(s)
	if err != nil {
		return err
	}

	if fi.IsDir() {
		s, err = filepath.Abs(s)
		if err != nil {
			return err
		}
		p.includePath = append(p.includePath, s)
		return nil
	}

	return fmt.Errorf("not a directory %q", s)
}

func (p *defaultParser) Parse(fileName string, r io.Reader, opts ...participle.ParseOption) (*script.Script, error) {
	return p.init(p.parser.Parse(fileName, r, opts...))
}

func (p *defaultParser) ParseBytes(fileName string, b []byte, opts ...participle.ParseOption) (*script.Script, error) {
	return p.init(p.parser.ParseBytes(fileName, b, opts...))
}

func (p *defaultParser) ParseString(fileName, src string, opts ...participle.ParseOption) (*script.Script, error) {
	return p.init(p.parser.ParseString(fileName, src, opts...))
}

func (p *defaultParser) ParseFile(fileName string, opts ...participle.ParseOption) (*script.Script, error) {
	return p.init(p.parseFile(fileName, opts...))
}

func (p *defaultParser) parseFile(fileName string, opts ...participle.ParseOption) (*script.Script, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Note: Do not wrap with init() here as this function is also used for importing scripts!
	return p.parser.Parse(fileName, f, opts...)
}

func (p *defaultParser) includeTopDec(s *script.Script, s1 *script.Script) error {
	for _, inc := range s1.Include {
		for _, path := range inc.Path {
			if err := p.include(s, s1.Pos, p.includePath, path); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *defaultParser) include(s *script.Script, pos lexer.Position, paths []string, file string) error {
	if s.Includes == nil {
		s.Includes = make(map[string]interface{})
	}

	// Locate the file within paths
	src, err := p.findFile(paths, file)
	if err != nil {
		return errors.Error(pos, err)
	}

	// To prevent an infinite loop, if we have already included a file, then don't include it
	if _, exists := s.Includes[src]; exists {
		return nil
	}
	s.Includes[src] = true

	s1, err := p.parseFile(src)
	if err != nil {
		return err
	}

	// Add any function definitions but do not include main()
	for _, td := range s1.FunDec {
		if td.Name != "main" {
			s.FunDec = append(s.FunDec, td)
		}
	}

	// Handle any includes in this file
	err = p.includeTopDec(s, s1)
	if err != nil {
		return err
	}

	return nil
}

func (p *defaultParser) findFile(paths []string, file string) (string, error) {
	if file == "" {
		return "", fmt.Errorf("include cannot be %q", file)
	}

	for _, path := range paths {
		fp := filepath.Join(path, file)
		fi, err := os.Stat(fp)
		if err != nil && !os.IsNotExist(err) {
			return "", err
		} else if err == nil && !fi.IsDir() {
			return filepath.Abs(fp)
		}
	}

	return "", os.ErrNotExist
}
