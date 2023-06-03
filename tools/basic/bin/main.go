package main

import (
	"github.com/peter-mount/go-basic/tools/basic"
	"github.com/peter-mount/go-kernel/v2"
	"log"
)

func main() {
	err := kernel.Launch(
		&basic.Basic{},
	)
	if err != nil {
		log.Fatal(err)
	}
}
