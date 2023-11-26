package handler

import (
	"context"
	"errors"
)

var ErrFiltered = errors.New("filtered")

type FilteringHandler[I any] struct {
	isOkay func(ctx context.Context, input I) bool
}

func NewFilteringHandler[I any](isOkay func(ctx context.Context, input I) bool) Handler[I, I] {
	return &FilteringHandler[I]{
		isOkay: isOkay,
	}
}

func (f *FilteringHandler[I]) Handle(ctx context.Context, input <-chan I) (context.Context, <-chan I) {
	select {
	case <-ctx.Done():
		return ctx, input
	case i := <-input:
		newCtx, cancel := context.WithCancelCause(ctx)
		o := make(chan I)
		go func() {
			defer close(o)
			o <- i
		}()
		if !f.isOkay(ctx, i) {
			cancel(ErrFiltered)
		}
		return newCtx, o
	}
}
