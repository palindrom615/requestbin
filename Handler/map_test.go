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

	// act
	ctx, res := handlerTest(add3Handler, 2)

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

	// act
	ctx, res := handlerTest(add3Handler, 2)

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
