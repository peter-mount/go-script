package io

import "github.com/peter-mount/go-script/packages"

func init() {
	packages.Register("os", &OS{})
}
