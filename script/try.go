package script

import (
	"context"
	"github.com/alecthomas/participle/v2/lexer"
)

// Try is Java's try {} catch {} finally{} construct
type Try struct {
	Pos lexer.Position

	Body       *Statement `parser:"'try' @@"`        // body
	CatchIdent string     `parser:"('catch' @Ident"` // catch var
	Catch      *Statement `parser:" @@)?"`           // catch block
	Finally    *Statement `parser:"('finally' @@)?"` // finally block
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
