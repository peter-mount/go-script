package errors

import (
	"errors"
	"fmt"
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	// error used to implement break
	breakError = errors.New("break")
	// error used to implement continue
	continueError = errors.New("continue")

	// VisitorExit is an error which will terminate the Visitor.
	// This is the same as any error occurring within a Visitor except that the final error
	// returned from specific handlers will become nil.
	VisitorExit = errors.New("visitor exit")

	// VisitorStop is an error which causes the current step in a Visitor to stop processing.
	// It's used to enable a Visitor to handle all processing of a node within itself rather
	// than the Visitor proceeding to any child nodes of that node.
	VisitorStop = errors.New("visitor stop")
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
	// Break and return dummy errors also are unchanged
	if err == nil || IsError(err) || IsBreak(err) || IsContinue(err) || IsReturn(err) || IsNoFieldErr(err) || IsVisitorStop(err) || IsVisitorExit(err) {
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

func NoField(pos lexer.Position, v interface{}, n string) error {
	return &NoFieldError{
		msg: fmt.Sprintf("%s %T has no field %q", pos.String(), v, n),
		v:   v,
		n:   n,
	}
}

type NoFieldError struct {
	msg string
	v   interface{}
	n   string
}

func (e *NoFieldError) Error() string {
	return e.msg
}

func (e *NoFieldError) Value() interface{} {
	return e.v
}

func (e *NoFieldError) Name() string {
	return e.n
}

func IsNoFieldErr(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*NoFieldError)
	return ok
}

func GetNoFieldErr(err error) (*NoFieldError, bool) {
	if err == nil {
		return nil, false
	}
	e, ok := err.(*NoFieldError)
	return e, ok
}

func Break() error {
	return breakError
}

// IsBreak returns true if err is from a break instruction being invoked.
func IsBreak(err error) bool { return err == breakError }

func Continue() error {
	return continueError
}

// IsContinue returns true if err is from a continue instruction being invoked.
func IsContinue(err error) bool { return err == continueError }

func IsReturn(err error) bool {
	_, ok := err.(*ReturnError)
	return ok
}

type ReturnError struct {
	value interface{}
}

func (r *ReturnError) Error() string {
	return "return"
}

func (r *ReturnError) Value() interface{} {
	return r.value
}

func NewReturn(v interface{}) error {
	return &ReturnError{value: v}
}

// IsVisitorStop returns true if err is VisitorStop
func IsVisitorStop(err error) bool {
	return err != nil && errors.Is(err, VisitorStop)
}

// IsVisitorExit returns true if err is VisitorExit
func IsVisitorExit(err error) bool {
	return err != nil && errors.Is(err, VisitorExit)
}
