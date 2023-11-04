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

func (h *ComposeHandler) Handle(ctx context.Context, input <-chan interface{}) (context.Context, <-chan interface{}, error) {
	out := make(chan interface{})
	lastOutput := input
	var err error
	for _, handler := range h.handlers {
		ctx, lastOutput, err = handler.Handle(ctx, lastOutput)
		if err != nil {
			return ctx, lastOutput, err
		}
	}
	out <- <-input
	return ctx, out, nil
}
