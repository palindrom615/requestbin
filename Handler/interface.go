package handler

import "context"

type Handler[I any, O any] interface {
	Handle(ctx context.Context, input <-chan I) (context.Context, <-chan O, error)
}

type IdentityHandler[I any] struct{}

func NewIdentityHandler[I any]() *IdentityHandler[I] {
	return &IdentityHandler[I]{}
}

func (h *IdentityHandler[I]) Handle(ctx context.Context, input <-chan I) (context.Context, <-chan I, error) {
	return ctx, input, nil
}
