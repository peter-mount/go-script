package executor

import (
	"errors"
	"fmt"
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	// Dummy errors for handling break and return statements
	breakError = errors.New("break")
	//returnError = errors.New("return")
)

type posError struct {
	msg string
}

func (e posError) Error() string {
	return e.msg
}

// Errorf returns an error containing the lexer.Position and the formatted message.
// IsError with this error will return true.
func Errorf(pos lexer.Position, f string, a ...interface{}) error {
	return &posError{msg: fmt.Sprintf(pos.String()+" "+f, a...)}
}

// Error wraps an error with the lexer.Position.
// If IsError(err) returns true then that error is returned.
// If err is nil then nil is returned.
// IsError will return true for a non-nil result from this function.
func Error(pos lexer.Position, err error) error {
	// If err is a PosError then return it as it has the position already.
	// Also, if err is nil then return nil, so we can use it as a catch-all
	if err == nil || IsError(err) {
		return err
	}
	return Errorf(pos, err.Error())
}

// IsError returns true if the error is from Errorf or Error functions.
func IsError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*posError)
	return ok
}

// IsBreak returns true if err is from a break instruction being invoked.
func IsBreak(err error) bool { return err == breakError }

func IsReturn(err error) bool {
	_, ok := err.(*returnError)
	return ok
}

type returnError struct {
	Value interface{}
}

func (r *returnError) Error() string {
	return "return"
}
