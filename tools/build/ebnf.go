package build

import (
	"github.com/peter-mount/go-build/core"
	"github.com/peter-mount/go-build/util/makefile"
	"github.com/peter-mount/go-build/util/makefile/target"
	"github.com/peter-mount/go-build/util/meta"
	"github.com/peter-mount/go-script/parser"
	"html"
	"os"
	"path/filepath"
	"strings"
)

// EBNF generates EBNF and railroad diagrams of the parser
type EBNF struct {
	Encoder *core.Encoder `kernel:"inject"`
	Build   *core.Build   `kernel:"inject"`
	Ebnf    *string       `kernel:"flag,ebnf,Output Parser EBNF"`
}

func (s *EBNF) Start() error {
	s.Build.Documentation(0, s.documentation)

	if *s.Ebnf != "" {
		return s.ebnf()
	}
	return nil
}

func (s *EBNF) documentation(_ makefile.Builder, target target.Builder, _ *meta.Meta) {

	srcDir := "documentation/go-script/parser"
	ebnf := filepath.Join(srcDir, "ebnf.txt")

	target.PhonyTarget(ebnf).
		MkDir(srcDir).
		Echo("GEN EBNF", ebnf).
		BuildTool("-ebnf", ebnf, "-d", *s.Encoder.Dest)
}

func (s *EBNF) ebnf() error {
	p := parser.New()

	ebnf := p.EBNF()

	err := os.MkdirAll(filepath.Dir(*s.Ebnf), 0755)

	// The raw ebnf file. Not needed in the final documentation
	// as we include it below, but needed for the railroad diagrams
	if err == nil {
		err = os.WriteFile(*s.Ebnf, []byte(ebnf), 0644)
	}

	// the html page forming part of the documentation
	if err == nil {
		err = os.WriteFile(
			filepath.Join(
				filepath.Dir(*s.Ebnf),
				strings.ReplaceAll(filepath.Base(*s.Ebnf), ".txt", ".html"),
			),
			[]byte(`---
editingNote: "This page is generated, any changes will be lost"
type: "manual"
title: "Extended Backus-Naur Form"
titleClass: section
linkTitle: "EBNF"
weight: 1
description: "Extended Backus-Naur Form of the scripting language"
---
<p>This is the Extended Backus-Naur Form (EBNF) description of the scripting language:</p>
<div class="sourceCode">`+html.EscapeString(ebnf)+`</div>
`),
			0644)
	}

	return err
}
