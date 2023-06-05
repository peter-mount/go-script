package executor

import "errors"

var (
	// Dummy errors for handling break and return statements
	breakError = errors.New("break")
	//returnError = errors.New("return")
)

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
