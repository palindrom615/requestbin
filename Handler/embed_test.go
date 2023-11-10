package handler_test

import (
	"context"
	"github.com/palindrom615/requestbin/handler"
	"testing"
)

func TestEmbedInputHandler_Handle(t *testing.T) {
	eih := handler.NewEmbedInputHandler(
		func(ctx context.Context, input interface{}) handler.CtxKey {
			return "key"
		},
	)
	i := make(chan interface{})
	go func() {
		i <- "value"
	}()
	ctx, _, err := eih.Handle(context.Background(), i)
	if err != nil {
		t.Error(err)
	}
	if ctx.Value(handler.CtxKey("key")) != "value" {
		t.Errorf("context value not set; expected 'value', got %v", ctx.Value("key"))
	}
}
