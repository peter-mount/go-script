package executor

import (
	"github.com/peter-mount/go-script/errors"
	"github.com/peter-mount/go-script/script"
)

// Statements executes a Statements block.
// Within this it runs with its own variable scope which is automatically closed when it completes
func (e *executor) Statements(statements *script.Statements) error {

	// Do nothing if it's an empty block
	if statements == nil || len(statements.Statements) == 0 {
		return nil
	}

	e.state.NewScope()
	defer e.state.EndScope()

	s := statements.Statements[0]
	for s != nil {
		if err := e.Statement(s); err != nil {
			return errors.Error(s.Pos, err)
		}
		s = s.Next
	}
	return nil
}

func (e *executor) Statement(statement *script.Statement) error {
	if statement == nil || statement.Empty {
		return nil
	}

	switch {
	case statement.Block != nil:
		return errors.Error(statement.Pos, e.Statements(statement.Block))

	case statement.Expression != nil:
		// Wrap visit to Expression, so we don't leak return values on the stack
		_, _, err := e.calculator.Calculate(func() error {
			return errors.Error(statement.Pos, e.Expression(statement.Expression))
		})
		return err

	case statement.For != nil:
		return errors.Error(statement.Pos, e.forStatement(statement.For))

	case statement.ForRange != nil:
		return errors.Error(statement.Pos, e.forRange(statement.ForRange))

	case statement.IfStmt != nil:
		return errors.Error(statement.Pos, e.ifStatement(statement.IfStmt))

	case statement.DoWhile != nil:
		return errors.Error(statement.Pos, e.doWhile(statement.DoWhile))

	case statement.Repeat != nil:
		return errors.Error(statement.Pos, e.repeatUntil(statement.Repeat))

	case statement.While != nil:
		return errors.Error(statement.Pos, e.while(statement.While))

	case statement.Return != nil:
		return errors.Error(statement.Pos, e.returnStatement(statement.Return))

	case statement.Break:
		return errors.Break()

	case statement.Continue:
		return errors.Continue()

	case statement.Switch != nil:
		return errors.Error(statement.Pos, e.switchStatement(statement.Switch))

	case statement.Try != nil:
		return errors.Error(statement.Pos, e.try(statement.Try))

	default:
		// This will fail if we add a new statement, but it's not yet implemented.
		return errors.Errorf(statement.Pos, "unimplemented statement reached")
	}
}
