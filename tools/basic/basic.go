package basic

import (
	"flag"
	"github.com/peter-mount/go-kernel/v2/log"
	common "github.com/peter-mount/go-script"
	"github.com/peter-mount/go-script/parser"
)

type Basic struct {
}

func (b *Basic) Run() error {
	for _, fileName := range flag.Args() {
		s, err := parser.ParseFile(fileName)
		if err != nil {
			return err
		}

		log.Println("Read", s)

		err = common.DebugScript(s)
		if err != nil {
			return err
		}
	}
	return nil
}
