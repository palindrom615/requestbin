package handler_test

import (
	"context"
	"testing"

	"github.com/palindrom615/requestbin/handler"
)

func TestEmbedInputHandler_Handle(t *testing.T) {
	eih := handler.NewEmbedCtxHandler(
		func(ctx context.Context, input interface{}) (handler.CtxKey, any) {
			return "key", "value"
		},
	)
	i := make(chan interface{})
	go func() {
		i <- "value"
	}()
	ctx, _ := eih.Handle(context.Background(), i)
	if ctx.Value(handler.CtxKey("key")) != "value" {
		t.Errorf("context value not set; expected 'value', got %v", ctx.Value("key"))
	}
}
