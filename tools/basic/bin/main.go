package main

import (
	"fmt"
	"github.com/peter-mount/go-basic/tools/basic"
	"github.com/peter-mount/go-kernel/v2"
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
