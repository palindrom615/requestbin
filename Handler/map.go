package handler

import (
	"context"
)

type MappingHandler[I any, O any] struct {
	mapper func(ctx context.Context, input I) (O, error)
}

func NewMappingHandler[I any, O any](mapper func(ctx context.Context, input I) (O, error)) Handler[I, O] {
	return &MappingHandler[I, O]{
		mapper: mapper,
	}
}

func (h *MappingHandler[I, O]) Handle(ctx context.Context, input <-chan I) (context.Context, <-chan O) {
	o := make(chan O)
	select {
	case <-ctx.Done():
		close(o)
		return ctx, o
	case i := <-input:
		res, err := h.mapper(ctx, i)
		newCtx, cancel := context.WithCancelCause(ctx)
		if err != nil {
			cancel(err)
		}
		go func() {
			defer close(o)
			o <- res
		}()
		return newCtx, o
	}
}
