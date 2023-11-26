package handler

import "context"

type ConsHandler[I, O1, O2 any] struct {
	handlerA Handler[I, O1]
	handlerB Handler[O1, O2]
}

func NewConsHandler[I, O1, O2 any](handlerA Handler[I, O1], handlerB Handler[O1, O2]) Handler[I, O2] {
	return &ConsHandler[I, O1, O2]{
		handlerA: handlerA,
		handlerB: handlerB,
	}
}

func (h *ConsHandler[I, O1, O2]) Handle(ctx context.Context, input <-chan I) (context.Context, <-chan O2) {
	return h.handlerB.Handle(h.handlerA.Handle(ctx, input))
}
