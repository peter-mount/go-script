package visitor

import "github.com/peter-mount/go-script/script"

func (v *visitor) VisitStatements(s *script.Statements) error {
	if s == nil {
		return nil
	}
	return v.visit(s.WithContext, func() error {
		if err := v.statements.Do(v.ctx); err != nil {
			return err
		}

		if !v.statementsNoNest {
			for _, e := range s.Statements {
				if err := v.VisitStatement(e); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (v *visitor) VisitStatement(s *script.Statement) error {
	if s == nil {
		return nil
	}
	return v.visit(s.WithContext, func() error {
		if err := v.statement.Do(v.ctx); err != nil {
			return err
		}

		switch {
		case s.Block != nil:
			return v.VisitStatements(s.Block)

		default:
			return nil
		}
	})
}

func (v *visitor) VisitFor(s *script.ForStmt) error {
	return v.visitTask(s.WithContext, v.forStatement)
}

func (v *visitor) VisitIf(s *script.IfStmt) error {
	return v.visitTask(s.WithContext, v.ifStatement)
}

func (v *visitor) VisitReturn(s *script.ReturnStmt) error {
	return v.visitTask(s.WithContext, v.returnStatement)
}

func (v *visitor) VisitWhile(s *script.WhileStmt) error {
	return v.visitTask(s.WithContext, v.whileStatement)
}
