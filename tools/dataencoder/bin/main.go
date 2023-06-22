package main

import (
	"fmt"
	"github.com/peter-mount/go-kernel/v2"
	_ "github.com/peter-mount/go-script/tools/dataencoder"
	"os"
)

func main() {
	if err := kernel.Launch(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
