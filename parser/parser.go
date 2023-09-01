package parser

import (
	"context"
	"fmt"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/script"
	"github.com/peter-mount/go-script/visitor"
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
	s, err := p.parser.Parse(fileName, r, opts...)
	if err != nil {
		return nil, err
	}
	return p.init(s)
}

func (p *defaultParser) ParseBytes(fileName string, b []byte, opts ...participle.ParseOption) (*script.Script, error) {
	s, err := p.parser.ParseBytes(fileName, b, opts...)
	if err != nil {
		return nil, err
	}
	return p.init(s)
}

func (p *defaultParser) ParseString(fileName, src string, opts ...participle.ParseOption) (*script.Script, error) {
	s, err := p.parser.ParseString(fileName, src, opts...)
	if err != nil {
		return nil, err
	}
	return p.init(s)
}

func (p *defaultParser) ParseFile(fileName string, opts ...participle.ParseOption) (*script.Script, error) {
	s, err := p.parseFile(fileName, opts...)
	if err != nil {
		return nil, err
	}
	return p.init(s)
}

func (p *defaultParser) parseFile(fileName string, opts ...participle.ParseOption) (*script.Script, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return p.parser.Parse(fileName, f, opts...)
}

func (p *defaultParser) includeTopDec(s *script.Script, tds []*script.TopDec) error {
	for _, td := range tds {
		if td.Include != nil {
			for _, path := range td.Include.Path {
				if err := p.include(s, td.Pos, p.includePath, path); err != nil {
					return err
				}
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
		return executor.Error(pos, err)
	}

	// To prevent an infinite loop, if we have already included a file, then dont include it
	if _, exists := s.Includes[src]; exists {
		return nil
	}
	s.Includes[src] = true

	s1, err := p.parseFile(src)
	if err != nil {
		return err
	}

	for _, td := range s1.TopDec {
		// Add any function definitions but do not include main()
		if td.FunDec != nil && td.FunDec.Name != "main" {
			s.TopDec = append(s.TopDec, td)
		}

		// Handle any includes in this file
		err = p.includeTopDec(s, s1.TopDec)
		if err != nil {
			return err
		}
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
		} else if !fi.IsDir() {
			return filepath.Abs(fp)
		}
	}

	return "", os.ErrNotExist
}

func (p *defaultParser) init(s *script.Script) (*script.Script, error) {
	err := p.includeTopDec(s, s.TopDec)
	if err != nil {
		return nil, err
	}

	err = visitor.New().
		Statements(p.initStatements).
		Statement(p.initStatement).
		If(p.initIf).
		While(p.initWhile).
		WithContext(context.Background()).
		VisitScript(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// initStatement sets Statements.Parent to point to the Statements
// instance containing this one.
//
// # It works by presuming that the immediate caller is usually Statement.Block
//
// The only time this should not be the case is when processing a
// function definition. As that's a top level object, there is no Statement
// in the context, hence Statements.Parent will be nil.
func (p *defaultParser) initStatements(ctx context.Context) error {
	// The Statements being visited
	statements := script.StatementsFromContext(ctx)
	if statements != nil {
		// The parent Statement, e.g. Statement.Block
		parent := script.StatementFromContext(ctx)
		if parent != nil {
			statements.Parent = parent.Parent
		}

		// Ensure Statement.Next is setup
		for i, s := range statements.Statements {
			if i > 0 {
				statements.Statements[i-1].Next = s
			}
		}
	}
	return nil
}

// initStatement sets Statement.Parent to the containing Statements instance.
//
// This works as Statement is only ever contained within a Statements instance
// so the parent is its parent.
func (p *defaultParser) initStatement(ctx context.Context) error {
	// The Statement being visited
	statement := script.StatementFromContext(ctx)
	if statement != nil {
		// This works because Statement is always inside a Statements instance
		parent := script.StatementsFromContext(ctx)
		if parent != nil {
			statement.Parent = parent
		}

		if statement.IfStmt != nil {
			if err := p.initIf(statement.IfStmt.WithContext(ctx)); err != nil {
				return err
			}
		}

		if statement.WhileStmt != nil {
			if err := p.initWhile(statement.WhileStmt.WithContext(ctx)); err != nil {
				return err
			}
		}

		if statement.ForRange != nil {
			if err := p.initForRange(statement.ForRange.WithContext(ctx)); err != nil {
				return err
			}
		}

		if statement.ForStmt != nil {
			if err := p.initFor(statement.ForStmt.WithContext(ctx)); err != nil {
				return err
			}
		}

		if statement.Try != nil {
			if err := p.initTry(statement.Try.WithContext(ctx)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *defaultParser) initIf(ctx context.Context) error {
	s := script.IfFromContext(ctx)
	v := visitor.FromContext(ctx)
	if err := v.VisitStatement(s.Body); err != nil {
		return err
	}
	if err := v.VisitStatement(s.Else); err != nil {
		return err
	}
	return nil
}

func (p *defaultParser) initWhile(ctx context.Context) error {
	s := script.WhileFromContext(ctx)
	v := visitor.FromContext(ctx)
	if err := v.VisitStatement(s.Body); err != nil {
		return err
	}
	return nil
}

func (p *defaultParser) initFor(ctx context.Context) error {
	s := script.ForFromContext(ctx)
	v := visitor.FromContext(ctx)
	if err := v.VisitStatement(s.Body); err != nil {
		return err
	}
	return nil
}

func (p *defaultParser) initForRange(ctx context.Context) error {
	s := script.ForRangeFromContext(ctx)
	v := visitor.FromContext(ctx)
	if err := v.VisitStatement(s.Body); err != nil {
		return err
	}
	return nil
}

func (p *defaultParser) initTry(ctx context.Context) error {
	s := script.TryFromContext(ctx)
	v := visitor.FromContext(ctx)

	// try-resources ensure only assignments and enforce declare mode
	// as those variables can only be accessed from within the body
	if s.Init != nil {
		for _, init := range s.Init.Resources {
			if init.Right != nil && init.Right.Right != nil {
				init.Right.Declare = true
			}
		}
	}

	if s.Body != nil {
		if err := v.VisitStatement(s.Body); err != nil {
			return err
		}
	}

	if s.Catch != nil {
		if err := v.VisitStatement(s.Catch.Statement); err != nil {
			return err
		}
	}

	if s.Finally != nil {
		if err := v.VisitStatement(s.Finally.Statement); err != nil {
			return err
		}
	}

	return nil
}
