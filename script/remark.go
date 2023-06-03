package script

import (
	"context"
	"github.com/alecthomas/participle/v2/lexer"
)

type Remark struct {
	Pos     lexer.Position
	Comment string `parser:"@Comment"`
}

func (s *Remark) Accept(v Visitor) error { return v.VisitRemark(s) }

func (s *Remark) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, remarkKey, s)
}

func RemarkFromContext(ctx context.Context) *Remark {
	return ctx.Value(remarkKey).(*Remark)
}
