package handler_test

import (
	"context"
	"testing"

	"github.com/palindrom615/requestbin/handler"
)

func TestEmbedInputHandler_Handle(t *testing.T) {
	// arrange
	eih := handler.NewEmbedCtxHandler(
		func(ctx context.Context, input interface{}) (handler.CtxKey, any) {
			return "key", "value"
		},
	)

	// act
	ctx, _ := handlerTest(eih, "value")

	// assert
	if ctx.Value(handler.CtxKey("key")) != "value" {
		t.Errorf("context value not set; expected 'value', got %v", ctx.Value("key"))
	}
}
