package main

import (
	"fmt"
	"github.com/peter-mount/go-kernel/v2"
	"github.com/peter-mount/go-script/tools/basic"
	"os"
)

func main() {
	err := kernel.Launch(
		&basic.Basic{},
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
