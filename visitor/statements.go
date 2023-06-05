package visitor

import "github.com/peter-mount/go-script/script"

func (v *visitor) VisitStatements(s *script.Statements) error {
	return v.visit(s.WithContext, func() error {
		if err := v.statements.Do(v.ctx); err != nil {
			return err
		}

		for _, e := range s.Statements {
			if err := v.VisitStatement(e); err != nil {
				return err
			}
		}
		return nil
	})
}

func (v *visitor) VisitStatement(s *script.Statement) error {
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
