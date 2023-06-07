package script

import (
	"context"
	"github.com/alecthomas/participle/v2/lexer"
)

// Try is Java's try {} catch {} finally{} construct.
//
// The Try.Init section contains variables pointing to resources.
// If they implement TryClosable then they will be closed when the
// body is exited.
//
// Try.Catch is called if an error occurs in the body. This is optiona.
//
// Try.Finally is called once the body & catch blocks have executed.
//
// Note: If a value in Try.Init implements TryClosable that interface is
// invoked before any catch or finally block. This follows the order
// in Java.
type Try struct {
	Pos lexer.Position

	Init       []*Expression `parser:"'try' [ '(' ( @@ (';' @@)* ) ')' ]"` // init block
	Body       *Statement    `parser:"@@"`                                 // body
	CatchIdent string        `parser:"('catch' @Ident"`                    // catch var
	Catch      *Statement    `parser:" @@)?"`                              // catch block
	Finally    *Statement    `parser:"('finally' @@)?"`                    // finally block
}

func (s *Try) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, tryKey, s)
}

func TryFromContext(ctx context.Context) *Try {
	v := ctx.Value(tryKey)
	if v != nil {
		return v.(*Try)
	}
	return nil
}

// TryClosable if implemented by a value defined in Try.Init will be called when the try block completes
type TryClosable interface {
	Close()
}
