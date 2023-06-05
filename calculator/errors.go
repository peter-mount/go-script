package calculator

import "errors"

var (
	stackEmpty = errors.New("stack empty")
)

func IsStackEmpty(err error) bool {
	return err == stackEmpty
}
