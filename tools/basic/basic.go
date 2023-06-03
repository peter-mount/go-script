package basic

import (
	"flag"
	common "github.com/peter-mount/go-basic"
	"github.com/peter-mount/go-basic/parser"
	"github.com/peter-mount/go-kernel/v2/log"
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
