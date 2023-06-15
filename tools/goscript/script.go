package goscript

import (
	"flag"
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/parser"
)

type Script struct {
}

func (b *Script) Run() error {
	p := parser.New()

	for _, fileName := range flag.Args() {
		s, err := p.ParseFile(fileName)
		if err != nil {
			return err
		}

		exec, err := executor.New(s)
		if err != nil {
			return err
		}

		err = exec.Run()
		if err != nil {
			return err
		}
	}

	return nil
}
