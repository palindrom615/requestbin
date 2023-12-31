package handler_test

import (
	"context"
	"testing"

	"github.com/palindrom615/requestbin/handler"
)

func handlerTest[I, O any](h handler.Handler[I, O], input I) (ctx context.Context, o <-chan O) {
	i := make(chan I)
	go func() {
		defer close(i)
		i <- input
	}()
	return h.Handle(context.Background(), i)
}

func TestIdentityHandler(t *testing.T) {
	// arrange
	h := handler.NewIdentityHandler[string]()

	// act
	_, o := handlerTest(h, "test")

	// assert
	if <-o != "test" {
		t.Errorf("IdentityHandler should return unchanged Context, got %v", o)
	}
}
