package debug

import (
	"context"
	"fmt"
	"github.com/peter-mount/go-script/state"
	"github.com/peter-mount/go-script/visitor"
	"strings"
)

func ListFunctions(s state.State) []string {
	var a []string
	var b []string

	v := visitor.New().
		FuncDec(func(_ context.Context) error {
			b = nil
			return nil
		}).
		WithContext(context.Background())

	funcs := s.GetFunctions()
	a = append(a, fmt.Sprintf("%d functions defined:", len(funcs)))

	retLen, funcLen := len("Return"), len("Function")
	for _, fName := range funcs {
		f, _ := s.GetFunction(fName)
		if funcLen < len(fName) {
			funcLen = len(fName)
		}
		if retLen < len(f.ReturnType) {
			retLen = len(f.ReturnType)
		}
	}

	format := fmt.Sprintf("%%%ds | %%%ds | %%s", retLen, funcLen)
	a = append(a, fmt.Sprintf(format, "Return", "Function", "Parameters"))
	for _, fName := range funcs {
		f, _ := s.GetFunction(fName)
		_ = v.VisitFuncDec(f)
		a = append(a, fmt.Sprintf(
			format,
			f.ReturnType,
			fName,
			strings.Join(b, ", ")),
		)
	}

	return a
}
