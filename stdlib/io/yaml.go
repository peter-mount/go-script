package io

import (
	"gopkg.in/yaml.v2"
	"io"
)

type YAML struct{}

func (j YAML) DecodeArray(r io.Reader) (any, error) {
	var m []any
	return j.decode(r, &m)
}

func (j YAML) DecodeMap(r io.Reader) (any, error) {
	m := make(map[string]interface{})
	return j.decode(r, &m)
}

func (_ YAML) decode(r io.Reader, v any) (any, error) {
	if err := yaml.NewDecoder(r).Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}

func (_ YAML) Encode(w io.Writer, v any) error {
	return yaml.NewEncoder(w).Encode(v)
}
