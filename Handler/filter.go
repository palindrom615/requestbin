package handler

import (
	"context"
	"errors"
)

var ErrFiltered = errors.New("filtered")

type FilteringHandler[I any] struct {
	filter func(ctx context.Context, input I) error
}

func NewFilteringHandler[I any](filter func(ctx context.Context, input I) error) Handler[I, I] {
	return &FilteringHandler[I]{
		filter: filter,
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
		err := f.filter(ctx, i)
		if err != nil {
			cancel(err)
		}
		return newCtx, o
	}
}
