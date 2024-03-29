package executor

import (
	"context"
	"github.com/peter-mount/go-kernel/v2/util/task"
	"github.com/peter-mount/go-script/errors"
	"github.com/peter-mount/go-script/script"
	"io"
)

func (e *executor) try(op *script.Try) (err error) {
	e.state.NewScope()
	defer e.state.EndScope()

	// Any panics get resolved to errors
	defer func() {
		if err1 := recover(); err1 != nil {
			err = errors.Errorf(op.Pos, "%v", err1)
		}
	}()

	if op.Finally != nil {
		defer func() {
			err1 := e.Statement(op.Finally.Statement)
			if err1 != nil {
				err = err1
			}
		}()
	}

	err = e.tryBody(op)
	if err != nil {
		if errors.IsReturn(err) {
			return err
		}

		err = errors.Error(op.Pos, err)

		// If catch then consume the error and pass it to the catch block
		if op.Catch != nil {
			// Set var unless "_" - always declared so always local
			if op.Catch.CatchIdent != "_" {
				e.state.Declare(op.Catch.CatchIdent)
				e.state.Set(op.Catch.CatchIdent, err.Error())
			}
			err = errors.Error(op.Pos, e.Statement(op.Catch.Statement))
		}
	}

	return
}

// tryBody runs any resources then the body.
// Note resources will be closed before any catch/finally blocks
func (e *executor) tryBody(op *script.Try) error {
	// Scope for resources & body
	e.state.NewScope()
	defer e.state.EndScope()

	// The deferable tasks to perform when we exit.
	//
	// we defer it here so that this task is always executed even if we don't get to
	// execute the action - e.g. creating a resource fails whilst building the list
	// means we still close any preceding resource in the resourceList
	var deferables task.Task
	defer func() {
		_ = deferables.Do(nil)
	}()

	// The action to perform
	var action task.Task

	// Configure any try-with-resources
	if op.Init != nil {
		for _, init := range op.Init.Resources {
			// Wrap visit to Expression, so we don't leak return values on the stack
			val, ok, err := e.calculator.Calculate(func() error {
				return errors.Error(init.Pos, e.Expression(init))
			})
			if err != nil {
				return err
			}

			if ok {
				// -----------------------------------------------------------------
				// NOTE: Always use deferables = task.Of(task).Defer(deferables) here so that
				// if the task fails, the rest of the deferables tasks still execute
				// -----------------------------------------------------------------
				if cl, ok := val.(CreateCloser); ok {
					if err := cl.Create(); err != nil {
						return errors.Error(init.Pos, err)
					}
				}

				// Common to io.Closer and CreateCloser
				if cl, ok := val.(io.Closer); ok {
					// add Close() from io.Closer to deferables
					deferables = task.Of(func(_ context.Context) error {
						return cl.Close()
					}).Defer(deferables)
				}
			}
		}
	}

	// Finally add the body to the action task
	if op.Body != nil {
		action = action.Then(func(_ context.Context) error {
			return errors.Error(op.Pos, e.Statement(op.Body))
		})
	}

	// Now execute the action. No need to add deferables as we have deferred its execution earlier
	return action.Do(nil)
}

// CreateCloser interface implemented by types that can be used as resources
type CreateCloser interface {
	io.Closer
	// Create is called when the resource is referenced before the statement is executed
	Create() error
}
