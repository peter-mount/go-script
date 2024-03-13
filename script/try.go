package script

import (
	"github.com/alecthomas/participle/v2/lexer"
)

// Try is based on Java's try {} catch {} finally{} construct.
//
// The Try.Init section contains variables pointing to resources.
// If they implement TryClosable then they will be closed when the
// body is exited.
//
// Try.Catch is called if an error occurs in the body. This is optional.
//
// Try.Finally is called once the body & catch blocks have executed.
//
// Note: If a value in Try.Init implements TryClosable that interface is
// invoked before any catch or finally block. This follows the order
// in Java.
type Try struct {
	Pos lexer.Position

	Init    *ResourceList `parser:"'try' @@?"` // init block
	Body    *Statement    `parser:"@@"`        // body
	Catch   *Catch        `parser:"@@?"`       // catch block
	Finally *Finally      `parser:"@@?"`       // finally block
}

type ResourceList struct {
	Resources []*Expression `parser:"'(' @@ (';' @@)* ')'"` // init block
}

type Catch struct {
	CatchIdent string     `parser:"'catch' '(' @Ident ')'"` // catch var
	Statement  *Statement `parser:" @@"`                    // catch block
}

type Finally struct {
	Statement *Statement `parser:"'finally' @@"` // finally block
}
