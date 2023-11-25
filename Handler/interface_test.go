package handler_test

import (
	"context"
	"testing"

	"github.com/palindrom615/requestbin/handler"
)

func TestIdentityHandler(t *testing.T) {
	// arrange
	h := handler.NewIdentityHandler[string]()

	// act
	ctx := context.Background()
	i := make(chan string)
	newCtx, o := h.Handle(ctx, i)
	go func() {
		i <- "test"
	}()

	// assert
	if ctx != newCtx {
		t.Errorf("IdentityHandler should return unchanged Context, got %v", ctx)
	}
	if <-o != "test" {
		t.Errorf("IdentityHandler should return unchanged Context, got %v", o)
	}
}
