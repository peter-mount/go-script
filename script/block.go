package script

import "context"

type Block struct {
	Statements []*Statement `parser:"@@*"`
}

func (s *Block) Accept(v Visitor) error { return v.VisitBlock(s) }

func (s *Block) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, blockKey, s)
}

func BlockFromContext(ctx context.Context) *Block {
	return ctx.Value(blockKey).(*Block)
}
