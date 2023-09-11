package executor

import (
	"context"
	"github.com/peter-mount/go-script/script"
)

// statements executes a Statements block.
// Within this it runs with its own variable scope which is automatically closed when it completes
func (e *executor) statements(ctx context.Context) error {
	statements := script.StatementsFromContext(ctx)

	// Do nothing if it's an empty block
	if statements == nil || len(statements.Statements) == 0 {
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
	// Do nothing for these.
	// Empty is the no-op statement.
	// Block has its own visitor
	case statement.Empty, statement.Block != nil:
		return nil

	case statement.Expression != nil:
		// Wrap visit to expression, so we don't leak return values on the stack
		_, _, err := e.calculator.Calculate(func(_ context.Context) error {
			return Error(statement.Pos, e.visitor.VisitExpression(statement.Expression))
		}, ctx)
		return err

	case statement.For != nil:
		return Error(statement.Pos, e.visitor.VisitFor(statement.For))

	case statement.ForRange != nil:
		return Error(statement.Pos, e.visitor.VisitForRange(statement.ForRange))

	case statement.IfStmt != nil:
		return Error(statement.Pos, e.visitor.VisitIf(statement.IfStmt))

	case statement.DoWhile != nil:
		return Error(statement.Pos, e.visitor.VisitDoWhile(statement.DoWhile))

	case statement.Repeat != nil:
		return Error(statement.Pos, e.visitor.VisitRepeat(statement.Repeat))

	case statement.While != nil:
		return Error(statement.Pos, e.visitor.VisitWhile(statement.While))

	case statement.Return != nil:
		return Error(statement.Pos, e.visitor.VisitReturn(statement.Return))

	case statement.Break:
		return Break()

	case statement.Continue:
		return Continue()

	case statement.Try != nil:
		return Error(statement.Pos, e.visitor.VisitTry(statement.Try))

	default:
		// This will fail if we add a new statement, but it's not yet implemented.
		return Errorf(statement.Pos, "unimplemented statement reached")
	}
}
