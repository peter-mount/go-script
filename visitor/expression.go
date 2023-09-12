package visitor

import "github.com/peter-mount/go-script/script"

func (v *visitor) VisitTernary(s *script.Ternary) error {
	return v.visitTask(s.WithContext, v.ternary)
}

func (v *visitor) VisitLogic(s *script.Logic) error {
	return v.visitTask(s.WithContext, v.logic)
}

func (v *visitor) VisitAddition(s *script.Addition) error {
	return v.visitTask(s.WithContext, v.addition)
}

func (v *visitor) VisitAssignment(s *script.Assignment) error {
	return v.visitTask(s.WithContext, v.assignment)
}

func (v *visitor) VisitComparison(s *script.Comparison) error {
	return v.visitTask(s.WithContext, v.comparison)
}

func (v *visitor) VisitEquality(s *script.Equality) error {
	return v.visitTask(s.WithContext, v.equality)
}

func (v *visitor) VisitExpression(s *script.Expression) error {
	return v.visit(s.WithContext, func() error {
		if err := v.expression.Do(v.ctx); err != nil {
			return err
		}

		if !v.expressionNoNest && s.Right != nil {
			return v.VisitAssignment(s.Right)
		}

		return nil
	})
}

func (v *visitor) VisitMultiplication(s *script.Multiplication) error {
	return v.visitTask(s.WithContext, v.multiplication)
}

func (v *visitor) VisitPrimary(s *script.Primary) error {
	return v.visitTask(s.WithContext, v.primary)
}

func (v *visitor) VisitUnary(s *script.Unary) error {
	return v.visitTask(s.WithContext, v.unary)
}
