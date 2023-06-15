package io

import "os"

type OS struct{}

func (_ OS) Create(n string) (*os.File, error) {
	return os.Create(n)
}

func (_ OS) Open(n string) (*os.File, error) {
	return os.Open(n)
}
