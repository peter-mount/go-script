package goscript

import (
	"errors"
	"flag"
	"github.com/peter-mount/go-build/application"
	"github.com/peter-mount/go-build/version"
	"github.com/peter-mount/go-kernel/v2/log"
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/parser"
	"os"
)

type Script struct {
}

func (b *Script) Run() error {
	if log.IsVerbose() {
		log.Println(version.Version)
	}

	p := parser.New()

	// if ../include exists then add it to the path
	if err := p.IncludePath(application.FileName(application.STATIC, "include")); err != nil && !os.IsNotExist(err) {
		return err
	}

	if err := p.IncludePath("."); err != nil && !os.IsNotExist(err) {
		return err
	}

	args := flag.Args()

	if len(args) == 0 {
		return errors.New("no scripts provided")
	}

	for _, fileName := range args {
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
