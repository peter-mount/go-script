package visitor

import (
	"github.com/peter-mount/go-script/script"
)

func (v *visitor) VisitFuncDec(s *script.FuncDec) error {
	return v.visit(s.WithContext, func() error {
		if err := v.funcDec.Do(v.ctx); err != nil {
			return err
		}

		if s.FunBody != nil {
			if err := v.VisitStatements(s.FunBody); err != nil {
				return err
			}
		}

		return nil
	})
}

func (v *visitor) VisitCallFunc(s *script.CallFunc) error {
	return v.visitTask(s.WithContext, v.callFunc)
}
