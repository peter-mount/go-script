package main

import (
	"github.com/peter-mount/go-basic/tools/dataencoder"
	"github.com/peter-mount/go-kernel/v2"
	"log"
)

func main() {
	err := kernel.Launch(
		&dataencoder.Build{},
	)
	if err != nil {
		log.Fatal(err)
	}
}
