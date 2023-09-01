package main

import (
	"fmt"
	_ "github.com/peter-mount/go-build/core"
	"github.com/peter-mount/go-kernel/v2"
	"github.com/peter-mount/go-script/tools/build"
	"os"
)

func main() {
	if err := kernel.Launch(
		&build.EBNF{},
		&build.Railroad{},
	); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
