package visitor

import (
	"github.com/peter-mount/go-script/script"
)

// visitTopDec handles TopDec and VarDec delegating to funcDec, arrayDec or scalarDec
func (v *visitor) visitTopDec(topDec *script.TopDec) error {
	switch {
	case topDec.FunDec != nil:
		return v.VisitFuncDec(topDec.FunDec)
	case topDec.VarDec != nil && topDec.VarDec.ArrayDec != nil:
		return v.visitTask(topDec.VarDec.ArrayDec.WithContext, v.arrayDec)
	case topDec.VarDec != nil && topDec.VarDec.ScalarDec != nil:
		return v.visitTask(topDec.VarDec.ScalarDec.WithContext, v.scalarDec)
	default:
	}
	return nil
}

func (v *visitor) VisitFuncDec(s *script.FuncDec) error {
	return v.visit(s.WithContext, func() error {
		if err := v.funcDec.Do(v.ctx); err != nil {
			return err
		}

		for _, param := range s.Parameters {
			if err := v.VisitParameter(param); err != nil {
				return err
			}
		}

		if s.FunBody != nil {
			// TODO visit VarDec here
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

func (v *visitor) VisitParameter(p *script.Parameter) error {
	switch {
	case p.Scalar != nil:
		return v.visitTask(p.Scalar.WithContext, v.scalarParamDec)
	case p.Array != nil:
		return v.visitTask(p.Array.WithContext, v.arrayParamDec)
	default:
		return nil
	}
}
