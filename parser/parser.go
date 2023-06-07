package parser

import (
	"context"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/peter-mount/go-script/script"
	"github.com/peter-mount/go-script/visitor"
	"io"
	"os"
)

type Parser interface {
	Parse(fileName string, r io.Reader, opts ...participle.ParseOption) (*script.Script, error)
	ParseBytes(fileName string, b []byte, opts ...participle.ParseOption) (*script.Script, error)
	ParseString(fileName, src string, opts ...participle.ParseOption) (*script.Script, error)
	ParseFile(fileName string, opts ...participle.ParseOption) (*script.Script, error)
}

type defaultParser struct {
	lexer  *lexer.StatefulDefinition
	parser *participle.Parser[script.Script]
}

func New() Parser {
	return &defaultParser{
		lexer:  scriptLexer,
		parser: scriptParser,
	}
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
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	s, err := p.parser.Parse(fileName, f, opts...)
	if err != nil {
		return nil, err
	}
	return p.init(s)
}

func (p *defaultParser) init(s *script.Script) (*script.Script, error) {
	err := visitor.New().
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

		if statement.ForStmt != nil {
			if err := p.initFor(statement.ForStmt.WithContext(ctx)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *defaultParser) initIf(ctx context.Context) error {
	s := script.IfFromContext(ctx)
	if s != nil {
		v := visitor.FromContext(ctx)
		if err := v.VisitStatement(s.Body); err != nil {
			return err
		}
		if err := v.VisitStatement(s.Else); err != nil {
			return err
		}
	}
	return nil
}

func (p *defaultParser) initWhile(ctx context.Context) error {
	s := script.WhileFromContext(ctx)
	if s != nil {
		v := visitor.FromContext(ctx)
		if err := v.VisitStatement(s.Body); err != nil {
			return err
		}
	}
	return nil
}

func (p *defaultParser) initFor(ctx context.Context) error {
	s := script.ForFromContext(ctx)
	if s != nil {
		v := visitor.FromContext(ctx)
		if err := v.VisitStatement(s.Body); err != nil {
			return err
		}
	}
	return nil
}
