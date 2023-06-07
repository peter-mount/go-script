package executor

import (
	"context"
	"github.com/peter-mount/go-script/script"
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
			err1 := e.visitor.VisitStatement(op.Finally)
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
			if op.CatchIdent != "_" {
				e.state.Declare(op.CatchIdent)
				e.state.Set(op.CatchIdent, err.Error())
			}
			err = Error(op.Pos, e.visitor.VisitStatement(op.Catch))
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

	// Configure andy try-with-resources
	for _, init := range op.Init {
		// Wrap visit to expression, so we don't leak return values on the stack
		val, ok, err1 := e.calculator.Calculate(func(_ context.Context) error {
			return Error(init.Pos, e.visitor.VisitExpression(init))
		}, ctx)
		if err1 != nil {
			return err1
		}

		// If implements TryClosable then defer it
		if ok {
			if cl, ok := val.(script.TryClosable); ok {
				// IDE will show this as a possible resource leak due to
				// defer being inside a for loop but in this instance
				// we actually want this
				defer cl.Close()
			}
		}
	}

	if op.Body != nil {
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
	}

	return nil
}
