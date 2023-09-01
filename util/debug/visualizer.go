package debug

import (
	"github.com/peter-mount/go-script/script"
	"strconv"
	"strings"
)

type visualizer struct {
	scr []string
}

// Visualize converts a script.Script into HTML so you can visualize what the parsed
// script looks like.
func Visualize(s *script.Script) string {
	v := visualizer{}
	v.append(`<html><body><style>
.title {font-weight: bold;}
.subtitle {font-weight: italic;}
.warning {color:red;}
.cell {
display: inline-block;
border: 1px solid black;
margin: 1px;
padding: 1px;
text-align: center;
vertical-align: top;
}
</style>`)

	v.Script(s)

	return strings.Join(append(v.scr, "</body></html>"), "")
}

func (v *visualizer) append(s ...string) {
	v.scr = append(v.scr, s...)
}

func (v *visualizer) cellStart(title, subTitle string) {
	v.append(`<div class="cell">`)
	if title != "" {
		v.append(`<div class="title">`, title, "</div>")
	}
	if subTitle != "" {
		v.append(`<div class="subtitle">`, subTitle, "</div>")
	}
}

func (v *visualizer) cellEnd() {
	v.append("</div>")
}

func (v *visualizer) Script(s *script.Script) {
	v.cellStart("script", "")
	for _, fd := range s.FunDec {
		v.FuncDec(fd)
	}
	v.cellEnd()
}

func (v *visualizer) FuncDec(s *script.FuncDec) {
	var p []string
	for _, param := range s.Parameters {
		p = append(p, param)
	}
	v.cellStart("FuncDec "+s.Name, strings.Join(p, ", "))
	if s.FunBody != nil && s.FunBody.Statements != nil {
		v.Statements(s.FunBody)
	}
	v.cellEnd()
}

func (v *visualizer) Statements(s *script.Statements) {
	v.cellStart("Statements", "")
	for _, stat := range s.Statements {
		v.Statement(stat)
	}
	v.cellEnd()
}

func (v *visualizer) Statement(s *script.Statement) {
	v.cellStart("Statement", "")

	switch {
	case s.Empty:
	case s.Expression != nil:
		v.Expression(s.Expression)
	default:
		v.append(`<div class="warning">Unsupported Statement</div>`)
	}

	v.cellEnd()
}

func (v *visualizer) Expression(s *script.Expression) {
	v.cellStart("Expression", "")
	v.Assignment(s.Right)
	v.cellEnd()
}

func (v *visualizer) Assignment(s *script.Assignment) {
	v.cellStart("Assignment", s.Op)
	if s.Left != nil {
		v.Logic(s.Left)
	}
	if s.Right != nil {
		v.Equality(s.Right)
	}
	v.cellEnd()
}

func (v *visualizer) Logic(s *script.Logic) {
	v.cellStart("Logic", s.Op)
	if s.Left != nil {
		v.Equality(s.Left)
	}
	if s.Right != nil {
		v.Logic(s.Right)
	}
	v.cellEnd()
}

func (v *visualizer) Equality(s *script.Equality) {
	v.cellStart("Equality", s.Op)
	if s.Left != nil {
		v.Comparison(s.Left)
	}
	if s.Right != nil {
		v.Equality(s.Right)
	}
	v.cellEnd()
}

func (v *visualizer) Comparison(s *script.Comparison) {
	v.cellStart("Comparison", s.Op)
	if s.Left != nil {
		v.Addition(s.Left)
	}
	if s.Right != nil {
		v.Comparison(s.Right)
	}
	v.cellEnd()
}

func (v *visualizer) Addition(s *script.Addition) {
	v.cellStart("Addition", s.Op)
	if s.Left != nil {
		v.Multiplication(s.Left)
	}
	if s.Right != nil {
		v.Addition(s.Right)
	}
	v.cellEnd()
}

func (v *visualizer) Multiplication(s *script.Multiplication) {
	v.cellStart("Multiplication", s.Op)
	if s.Left != nil {
		v.Unary(s.Left)
	}
	if s.Right != nil {
		v.Multiplication(s.Right)
	}
	v.cellEnd()
}

func (v *visualizer) Unary(s *script.Unary) {
	v.cellStart("Unary", s.Op)
	if s.Left != nil {
		v.Unary(s.Left)
	}
	if s.Right != nil {
		v.Primary(s.Right)
	}
	v.cellEnd()
}

func (v *visualizer) Primary(s *script.Primary) {
	switch {
	case s.Null, s.Nil:
		v.cellStart("null", "")
	case s.True:
		v.cellStart("true", "")
	case s.False:
		v.cellStart("false", "")
	case s.Ident != "":
		v.cellStart("ident", s.Ident)
	case s.Integer != nil:
		v.cellStart("int ", strconv.Itoa(*s.Integer))
	case s.Float != nil:
		v.cellStart("float ", strconv.FormatFloat(*s.Float, 'f', 3, 64))
	case s.SubExpression != nil:
		v.cellStart("SubExpression", "()")
		v.Expression(s.SubExpression)
	default:
		v.cellStart("Primary", "???")
	}
	v.cellEnd()
}
