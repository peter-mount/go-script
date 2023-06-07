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

	if op.Body != nil {
		err = Error(op.Pos, e.visitor.VisitStatement(op.Body))
		if err != nil {
			if IsReturn(err) {
				return err
			}
			if IsBreak(err) {
				return nil
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
	}

	return
}
