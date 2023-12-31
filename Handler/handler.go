package handler

import (
	"context"
	"github.com/palindrom615/requestbin"
)

var logger = requestbin.GetLogger()

type Handler[I any, O any] interface {
	Handle(ctx context.Context, input <-chan I) (context.Context, <-chan O)
}

type IdentityHandler[I any] struct{}

func NewIdentityHandler[I any]() Handler[I, I] {
	return &IdentityHandler[I]{}
}

func (h *IdentityHandler[I]) Handle(ctx context.Context, input <-chan I) (context.Context, <-chan I) {
	return ctx, input
}
