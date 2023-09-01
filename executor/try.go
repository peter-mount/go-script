package executor

import (
	"context"
	"github.com/peter-mount/go-kernel/v2/util/task"
	"github.com/peter-mount/go-script/script"
	"io"
)

func (e *executor) try(ctx context.Context) (err error) {
	op := script.TryFromContext(ctx)

	e.state.NewScope()
	defer e.state.EndScope()

	// Any panics get resolved to errors
	defer func() {
		if err1 := recover(); err1 != nil {
			err = Errorf(op.Pos, "%v", err1)
		}
	}()

	if op.Finally != nil {
		defer func() {
			err1 := e.visitor.VisitStatement(op.Finally.Statement)
			if err1 != nil {
				err = err1
			}
		}()
	}

	err = e.tryBody(op, ctx)
	if err != nil {
		if IsReturn(err) {
			return err
		}

		err = Error(op.Pos, err)

		// If catch then consume the error and pass it to the catch block
		if op.Catch != nil {
			// Set var unless "_" - always declared so always local
			if op.Catch.CatchIdent != "_" {
				e.state.Declare(op.Catch.CatchIdent)
				e.state.Set(op.Catch.CatchIdent, err.Error())
			}
			err = Error(op.Pos, e.visitor.VisitStatement(op.Catch.Statement))
		}
	}

	return
}

// tryBody runs any resources then the body.
// Note resources will be closed before any catch/finally blocks
func (e *executor) tryBody(op *script.Try, ctx context.Context) error {
	// Scope for resources & body
	e.state.NewScope()
	defer e.state.EndScope()

	var action task.Task

	if op.Body != nil {
		action = action.Then(func(_ context.Context) error {
			err := Error(op.Pos, e.visitor.VisitStatement(op.Body))
			if err != nil {
				if IsReturn(err) {
					return err
				}
				if IsBreak(err) {
					return nil
				}

				return Error(op.Pos, err)
			}
			return nil
		})
	}

	// Configure andy try-with-resources
	if op.Init != nil {
		for _, init := range op.Init.Resources {
			// Wrap visit to expression, so we don't leak return values on the stack
			val, ok, err1 := e.calculator.Calculate(func(_ context.Context) error {
				return Error(init.Pos, e.visitor.VisitExpression(init))
			}, ctx)
			if err1 != nil {
				return err1
			}

			if ok {
				if cl, ok := val.(CreateCloser); ok {
					// CreateCloser will allow us to have a cope created by Create but closed when the try block completes
					err1 := cl.Create()
					if err1 != nil {
						return err1
					}

					action = action.Defer(func(_ context.Context) error {
						return cl.Close()
					})
				} else if cl, ok := val.(io.Closer); ok {
					action = action.Defer(func(_ context.Context) error {
						return cl.Close()
					})
				}
			}
		}
	}

	return action.Do(ctx)
}

type CreateCloser interface {
	io.Closer
	Create() error
}
