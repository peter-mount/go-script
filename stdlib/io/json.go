package io

import (
	"encoding/json"
	"io"
)

type JSON struct{}

func (j JSON) DecodeArray(r io.Reader) (any, error) {
	var m []any
	return j.decode(r, &m)
}

func (j JSON) DecodeMap(r io.Reader) (any, error) {
	m := make(map[string]interface{})
	return j.decode(r, &m)
}

func (_ JSON) decode(r io.Reader, v any) (any, error) {
	if err := json.NewDecoder(r).Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}

func (_ JSON) Encode(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}
