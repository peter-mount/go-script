package main

import (
	"fmt"
	"github.com/peter-mount/go-kernel/v2"
	"github.com/peter-mount/go-script/tools/goscript"
	"os"
)

func main() {
	err := kernel.Launch(
		&goscript.Script{},
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
