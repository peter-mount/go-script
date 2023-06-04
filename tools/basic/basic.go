package basic

import (
	"flag"
	"fmt"
	"github.com/peter-mount/go-script/debug"
	"github.com/peter-mount/go-script/parser"
	"github.com/peter-mount/go-script/state"
	"strings"
)

type Basic struct {
}

func (b *Basic) Run() error {
	for _, fileName := range flag.Args() {
		s, err := parser.ParseFile(fileName)
		if err != nil {
			return err
		}

		st, err := state.New(s)
		if err == nil {
			fmt.Println(strings.Join(debug.ListFunctions(st), "\n"))
			err = debug.DebugScript(s)
		}

		if err != nil {
			return err
		}
	}
	return nil
}
