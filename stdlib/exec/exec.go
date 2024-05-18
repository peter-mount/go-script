// Package exec provides access to Go's exec package
package exec

import (
	"github.com/peter-mount/go-script/packages"
	"os"
	"os/exec"
)

func init() {
	packages.Register("exec", &Exec{})
}

type Exec struct{}

// Run will execute the provided command, returning an error if the command fails.
// Stdin, Stdout and Stderr will be those of the script.
func (_ Exec) Run(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
