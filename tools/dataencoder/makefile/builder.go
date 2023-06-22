package makefile

import (
	"fmt"
	"strings"
)

type Builder interface {
	Comment(f string, a ...any) Builder
	Include(f string, a ...any) Builder
	Line(f string, a ...any) Builder
	Rule(string, ...string) Builder
	End() Builder
	Build() string
}

type builder struct {
	parent   *builder // parent builder
	key      string
	line     string
	comment  string
	command  string
	children []*builder
}

func New() Builder {
	return &builder{}
}

func (b *builder) End() Builder {
	if b.parent == nil {
		return b
	}
	return b.parent
}

func (b *builder) Comment(f string, a ...any) Builder {
	c := &builder{parent: b, comment: fmt.Sprintf(f, a...)}
	b.children = append(b.children, c)
	return b
}

func (b *builder) Include(f string, a ...any) Builder {
	c := &builder{
		parent:  b,
		command: fmt.Sprintf("include "+f, a...),
	}
	b.children = append(b.children, c)
	return b
}

func (b *builder) Line(f string, a ...any) Builder {
	c := &builder{parent: b, line: fmt.Sprintf(f, a...)}
	b.children = append(b.children, c)
	return b
}

func (b *builder) Rule(rule string, targets ...string) Builder {
	c := &builder{
		parent: b,
		key:    rule,
		line:   strings.Join(targets, " "),
	}
	b.children = append(b.children, c)
	return b
}

func (b *builder) Build() string {
	return strings.Join(b.build([]string{}), "\n")
}

func (b *builder) build(a []string) []string {
	switch {
	case b.comment != "":
		return append(a, "# "+b.comment)

	case b.command != "":
		return append(a, b.command)

	case b.key != "":
		return append(a, "", b.key+": "+b.line)

	case b.line != "":
		return append(a, "\t"+b.line)

	default:
		for _, c := range b.children {
			a = c.build(a)
		}
	}
	return a
}
