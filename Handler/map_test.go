package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/palindrom615/requestbin/handler"
)

func TestMappingHandler_map_without_err(t *testing.T) {
	// arrange
	add3Handler := handler.NewMappingHandler(
		func(ctx context.Context, input int) (int, error) {
			return input + 3, nil
		},
	)
	i := make(chan int)
	go func() {
		defer close(i)
		i <- 2
	}()

	// act
	ctx, res := add3Handler.Handle(context.Background(), i)

	// assert
	select {
	case <-ctx.Done():
		t.Errorf("handler should not cancel")
	case r := <-res:
		if r != 5 {
			t.Errorf("result should be 5, got %d", r)
		}
	}
}

func TestMappingHandler_map_with_err(t *testing.T) {
	// arrange
	errThrown := errors.New("test error")
	add3Handler := handler.NewMappingHandler(
		func(ctx context.Context, input int) (int, error) {
			return 0, errThrown
		},
	)
	i := make(chan int)
	go func() {
		defer close(i)
		i <- 2
	}()

	// act
	ctx, res := add3Handler.Handle(context.Background(), i)

	// assert
	<-ctx.Done()
	if context.Cause(ctx) != errThrown {
		t.Errorf("error should be errThrown, got %v", context.Cause(ctx))
	}
	r := <-res
	if r != 0 {
		t.Errorf("handler should return 0, got %v", r)
	}
}
