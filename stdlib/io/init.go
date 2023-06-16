package io

import "github.com/peter-mount/go-script/packages"

func init() {
	packages.Register("io", newIO())
	packages.Register("json", &JSON{})
	packages.Register("os", &OS{})
	packages.Register("yaml", &YAML{})
}
