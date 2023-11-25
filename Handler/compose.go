package handler

import "context"

type ComposeHandler struct {
	handlers []Handler[interface{}, interface{}]
}

func NewComposeHandler(handlers ...Handler[interface{}, interface{}]) *ComposeHandler {
	return &ComposeHandler{
		handlers: handlers,
	}
}

func (h *ComposeHandler) Handle(ctx context.Context, input <-chan interface{}) (context.Context, <-chan interface{}) {
	lastOutput := input
	for _, handler := range h.handlers {
		ctx, lastOutput = handler.Handle(ctx, lastOutput)
	}
	return ctx, lastOutput
}
