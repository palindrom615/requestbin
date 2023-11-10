package handler

import "context"

type CtxKey string
type EmbedInputHandler[I any] struct {
	calcKey func(ctx context.Context, input I) CtxKey
}

func NewEmbedInputHandler[I any](calcKey func(ctx context.Context, input I) CtxKey) *EmbedInputHandler[I] {
	return &EmbedInputHandler[I]{
		calcKey: calcKey,
	}
}

func (h *EmbedInputHandler[I]) Handle(ctx context.Context, input <-chan I) (context.Context, <-chan I, error) {
	i := <-input
	out := make(chan I)
	go func() {
		out <- i
	}()
	key := h.calcKey(ctx, i)
	return context.WithValue(ctx, key, i), out, nil
}
