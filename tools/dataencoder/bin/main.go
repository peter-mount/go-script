package main

import (
	"fmt"
	"github.com/peter-mount/go-kernel/v2"
	"github.com/peter-mount/go-script/tools/dataencoder"
	"os"
)

func main() {
	err := kernel.Launch(
		&dataencoder.Build{},
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
