package handler

import "context"

type SequentialHandler struct {
	handlers []Handler[any, any]
}

func NewSequentialHandler(handlers ...Handler[any, any]) Handler[any, any] {
	return &SequentialHandler{
		handlers: handlers,
	}
}

func (h *SequentialHandler) Handle(ctx context.Context, input <-chan any) (context.Context, <-chan any) {
	lastOutput := input
	for _, handler := range h.handlers {
		ctx, lastOutput = handler.Handle(ctx, lastOutput)
	}
	return ctx, lastOutput
}
