package visitor

import "github.com/peter-mount/go-script/script"

func (v *visitor) VisitExpression(s *script.Expression) error {
	return v.visit(s.WithContext, func() error {
		return v.expression.Do(v.ctx)
	})
}
