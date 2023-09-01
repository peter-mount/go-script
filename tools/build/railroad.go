package build

import (
	"fmt"
	"github.com/alecthomas/participle/v2/ebnf"
	"github.com/alecthomas/repr"
	"github.com/peter-mount/go-build/core"
	"github.com/peter-mount/go-build/util/makefile/target"
	"github.com/peter-mount/go-build/util/meta"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Railroad creates railroad diagrams from EBMF files.
//
// This is based on the railroad command in participle
type Railroad struct {
	Encoder  *core.Encoder `kernel:"inject"`
	Build    *core.Build   `kernel:"inject"`
	Railroad *string       `kernel:"flag,railroad,Generate Railroad graph"`
	Output   *string       `kernel:"flag,railroad-out,Destination of Railroad graph"`
}

func (r *Railroad) Start() error {
	r.Build.Documentation(r.documentation)

	if *r.Railroad != "" {
		return r.run()
	}
	return nil
}

func (r *Railroad) documentation(target target.Builder, _ *meta.Meta) {

	// Write these to our documentation
	srcDir := "documentation/go-script/parser"
	ebnf := filepath.Join(srcDir, "ebnf.txt")

	output := filepath.Join(srcDir, "railroad.html")

	target.PhonyTarget(output, ebnf).
		MkDir(srcDir).
		Echo("RAILROAD", output).
		BuildTool(
			"-railroad", ebnf,
			"-railroad-out", output,
			"-d", *r.Encoder.Dest,
		)
}

const (
	mergeRefThreshold  = -1
	mergeSizeThreshold = 0
)

type production struct {
	*ebnf.Production
	refs int
	size int
}

func (r *Railroad) generate(productions map[string]*production, n ebnf.Node) (s string) {
	switch n := n.(type) {
	case *ebnf.EBNF:
		for _, p := range n.Productions {
			s += r.generate(productions, p)
		}

	case *ebnf.Production:
		if productions[n.Production].refs <= mergeRefThreshold {
			break
		}

		s += `    f.write('<h3 class="paragraph">` + n.Production + `</h3>\n')` + "\n"
		s += "    Diagram("
		s += r.generate(productions, n.Expression)
		s += ").writeSvg(f.write)\n"

	case *ebnf.Expression:
		s += "Choice(0, "
		for i, a := range n.Alternatives {
			if i > 0 {
				s += ", "
			}
			s += r.generate(productions, a)
		}
		s += ")"

	case *ebnf.SubExpression:
		s += r.generate(productions, n.Expr)
		if n.Lookahead != ebnf.LookaheadAssertionNone {
			s = fmt.Sprintf(`Group(%s, "?%c")`, s, n.Lookahead)
		}

	case *ebnf.Sequence:
		s += "Sequence("
		for i, t := range n.Terms {
			if i > 0 {
				s += ", "
			}
			s += r.generate(productions, t)
		}
		s += ")"

	case *ebnf.Term:
		switch n.Repetition {
		case "*":
			s += "ZeroOrMore("
		case "+":
			s += "OneOrMore("
		case "?":
			s += "Optional("
		}
		switch {
		case n.Name != "":
			p := productions[n.Name]
			if p.refs > mergeRefThreshold {
				//s += fmt.Sprintf("NonTerminal(%q, {href:\"#%s\"})", n.Name, n.Name)
				s += fmt.Sprintf("NonTerminal(%q)", n.Name)
			} else {
				s += r.generate(productions, p.Expression)
			}

		case n.Group != nil:
			s += r.generate(productions, n.Group)

		case n.Literal != "":
			s += fmt.Sprintf("Terminal(%s)", n.Literal)

		case n.Token != "":
			s += fmt.Sprintf("NonTerminal(%q)", n.Token)

		default:
			panic(repr.String(n))

		}
		if n.Repetition != "" {
			s += ")"
		}
		if n.Negation {
			s = fmt.Sprintf(`Group(%s, "~")`, s)
		}

	default:
		panic(repr.String(n))
	}
	return
}

func (r *Railroad) countProductions(productions map[string]*production, n ebnf.Node) (size int) {
	switch n := n.(type) {
	case *ebnf.EBNF:
		for _, p := range n.Productions {
			productions[p.Production] = &production{Production: p}
		}
		for _, p := range n.Productions {
			r.countProductions(productions, p)
		}
		for _, p := range n.Productions {
			if productions[p.Production].size <= mergeSizeThreshold {
				productions[p.Production].refs = mergeRefThreshold
			}
		}
	case *ebnf.Production:
		productions[n.Production].size = r.countProductions(productions, n.Expression)
	case *ebnf.Expression:
		for _, a := range n.Alternatives {
			size += r.countProductions(productions, a)
		}
	case *ebnf.SubExpression:
		size += r.countProductions(productions, n.Expr)
	case *ebnf.Sequence:
		for _, t := range n.Terms {
			size += r.countProductions(productions, t)
		}
	case *ebnf.Term:
		if n.Name != "" {
			productions[n.Name].refs++
			size++
		} else if n.Group != nil {
			size += r.countProductions(productions, n.Group)
		} else {
			size++
		}
	default:
		panic(repr.String(n))
	}
	return
}

func (r *Railroad) parse() (*ebnf.EBNF, error) {
	//f, err := os.Open(*r.Railroad)
	//if err != nil {
	//	return nil, err
	//}
	//defer f.Close()
	//return ebnf.Parse(f)
	b, err := os.ReadFile(*r.Railroad)
	if err != nil {
		return nil, err
	}
	s := string(b)
	s = strings.ReplaceAll(s, ")+?)", ")+)")
	return ebnf.ParseString(s)
}

func (r *Railroad) run() error {

	ast, err := r.parse()
	if err != nil {
		return err
	}

	productions := map[string]*production{}
	r.countProductions(productions, ast)
	str := r.generate(productions, ast)

	// Prefix the generated python with some boilerplate stuff
	str = `# pip install git+https://github.com/[repo owner]/[repo]@[branch name]
# pip install railroad-diagrams

import sys
from railroad import Diagram, Choice, Group, Sequence, ZeroOrMore, OneOrMore, Optional, Terminal, NonTerminal

with open('` + *r.Output + `', 'w', encoding="utf-8") as f:
    f.write('---\n')
    f.write('editingNote: "This page is generated, any changes will be lost"\n')
    f.write('type: "manual"\n')
    f.write('title: "Rail Road diagrams"\n')
    f.write('titleClass: chapter\n')
    f.write('linkTitle: "Rail Road"\n')
    f.write('weight: 1\n')
    f.write('description: "Rail Road diagrams of the parser"\n')
    f.write('---\n')
` + str

	err = os.MkdirAll(*r.Encoder.Dest, 0755)
	if err != nil {
		return err
	}

	python := filepath.Join(*r.Encoder.Dest, "genrailroad.py")
	err = os.WriteFile(python, []byte(str), 0644)
	if err != nil {
		return err
	}

	return r.runPython(python)
}

func (r *Railroad) runPython(script string) error {
	cmd := exec.Command("python", script)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
