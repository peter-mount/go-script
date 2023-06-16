package fmt

import "github.com/peter-mount/go-script/packages"

func init() {
	packages.Register("fmt", &FMT{})
}

type FMT struct{}
