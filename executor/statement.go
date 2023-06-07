package executor

import (
	"context"
	"github.com/peter-mount/go-script/script"
)

// statements executes a Statements block.
// Within this it runs with its own variable scope which is automatically closed when it completes
func (e *executor) statements(ctx context.Context) error {
	statements := script.StatementsFromContext(ctx)

	if statements == nil {
		return nil
	}

	e.state.NewScope()
	defer e.state.EndScope()

	s := statements.Statements[0]
	for s != nil {
		if err := e.visitor.VisitStatement(s); err != nil {
			return Error(s.Pos, err)
		}
		s = s.Next
	}
	return nil
}

func (e *executor) statement(ctx context.Context) error {
	statement := script.StatementFromContext(ctx)

	switch {
	case statement.Expression != nil:
		// Wrap visit to expression, so we don't leak return values on the stack
		_, _, err := e.calculator.Calculate(func(_ context.Context) error {
			return Error(statement.Pos, e.visitor.VisitExpression(statement.Expression))
		}, ctx)
		return err

	case statement.ForStmt != nil:
		return Error(statement.Pos, e.visitor.VisitFor(statement.ForStmt))

	case statement.ForRange != nil:
		return Error(statement.Pos, e.visitor.VisitForRange(statement.ForRange))

	case statement.IfStmt != nil:
		return Error(statement.Pos, e.visitor.VisitIf(statement.IfStmt))

	case statement.WhileStmt != nil:
		return Error(statement.Pos, e.visitor.VisitWhile(statement.WhileStmt))

	case statement.ReturnStmt != nil:
		return Error(statement.Pos, e.visitor.VisitReturn(statement.ReturnStmt))

	case statement.Break != nil:
		return Break()

	case statement.Try != nil:
		return Error(statement.Pos, e.visitor.VisitTry(statement.Try))
	}

	return nil
}
