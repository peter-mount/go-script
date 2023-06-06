package visitor

import (
	"github.com/peter-mount/go-script/script"
)

// visitTopDec handles TopDec and VarDec delegating to funcDec, arrayDec or scalarDec
func (v *visitor) visitTopDec(topDec *script.TopDec) error {
	switch {
	case topDec.FunDec != nil:
		return v.VisitFuncDec(topDec.FunDec)
	default:
	}
	return nil
}

func (v *visitor) VisitFuncDec(s *script.FuncDec) error {
	return v.visit(s.WithContext, func() error {
		if err := v.funcDec.Do(v.ctx); err != nil {
			return err
		}

		if s.FunBody != nil {
			if err := v.VisitStatements(s.FunBody.Statements); err != nil {
				return err
			}
		}

		return nil
	})
}

func (v *visitor) VisitFuncBody(s *script.FuncBody) error {
	return v.visitTask(s.WithContext, v.funcBody)
}

func (v *visitor) VisitCallFunc(s *script.CallFunc) error {
	return v.visitTask(s.WithContext, v.callFunc)
}
