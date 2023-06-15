package main

import (
	"fmt"
	"github.com/peter-mount/go-kernel/v2"
	_ "github.com/peter-mount/go-script/stdlib"
	_ "github.com/peter-mount/go-script/stdlib/math"
	"github.com/peter-mount/go-script/tools/goscript"
	"os"
)

func main() {
	if err := kernel.Launch(
		&goscript.Script{},
	); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
