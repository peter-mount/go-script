package goscript

import (
	"flag"
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/parser"
	"os"
	"path"
	"path/filepath"
)

type Script struct {
}

func (b *Script) Run() error {
	p := parser.New()

	// if ../include exists then add it to the path
	if err := p.IncludePath(path.Join(filepath.Dir(os.Args[0]), "../include")); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	if err := p.IncludePath("."); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

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
