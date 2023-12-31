package handler

import (
	"context"
)

type CtxKey any

type EmbedCtxHandler[I, V any] struct {
	getNewCtxVal func(ctx context.Context, input I) (CtxKey, V)
}

func NewEmbedCtxHandler[I, V any](getNewCtxVal func(ctx context.Context, input I) (CtxKey, V)) Handler[I, I] {
	return &EmbedCtxHandler[I, V]{
		getNewCtxVal: getNewCtxVal,
	}
}

func (h *EmbedCtxHandler[I, V]) Handle(ctx context.Context, input <-chan I) (context.Context, <-chan I) {
	select {
	case i := <-input:
		out := make(chan I)
		go func() {
			out <- i
		}()
		key, val := h.getNewCtxVal(ctx, i)
		logger.Debugw("new context value", "key", key, "value", val)
		return context.WithValue(ctx, key, val), out
	case <-ctx.Done():
		return ctx, input
	}
}
