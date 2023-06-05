package executor

import (
	"context"
	"fmt"
	"github.com/peter-mount/go-script/script"
)

// statements executes a Statements block.
// Within this it runs with its own variable scope which is automatically closed when it completes
func (e *executor) statements(ctx context.Context) error {
	statements := script.StatementsFromContext(ctx)
	fmt.Printf("%s %T\n", statements.Pos, statements)

	if statements == nil {
		panic("no statements?")
	}

	e.state.NewScope()
	defer e.state.EndScope()

	s := statements.Statements[0]
	for s != nil {
		if err := e.visitor.VisitStatement(s); err != nil {
			return err
		}
		s = s.Next
	}
	return nil
}

func (e *executor) statement(ctx context.Context) error {
	statement := script.StatementFromContext(ctx)
	fmt.Printf("%s %T\n", statement.Pos, statement)

	switch {
	case statement.Expression != nil:
		return e.visitor.VisitExpression(statement.Expression)

	}

	return nil
}
