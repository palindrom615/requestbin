package handler_test

import (
	"context"
	"testing"

	"github.com/palindrom615/requestbin/handler"
)

func isOkay(ctx context.Context, input string) bool {
	return input == "okay"
}
func TestFilteringHandler_HandleOkay(t *testing.T) {
	// arrange
	handler := handler.NewFilteringHandler(isOkay)

	// act
	newCtx, o := handlerTest(handler, "okay")

	// assert
	select {
	case s := <-o:
		t.Logf("Handler returned %s", s)
		if s != "okay" {
			t.Errorf("filtering handler should return \"okay\", got \"%s\"", s)
		}
	case <-newCtx.Done():
		t.Errorf("filtering handler should not be canceled: %s", newCtx.Err())
	}
}

func TestFilteringHandler_HandleNotokay(t *testing.T) {
	// arrange
	h := handler.NewFilteringHandler(isOkay)

	// act
	newCtx, o := handlerTest(h, "notOkay")

	// assert
	<-newCtx.Done()
	if context.Cause(newCtx) != handler.ErrFiltered {
		t.Errorf("FilteringHandler should be canceled because of handler.ErrFiltered, got %s", newCtx.Err())
	}
	s := <-o
	t.Logf("Handler returned %s", s)
	if s != "notOkay" {
		t.Errorf("FilteringHandler should return \"notOkay\", got \"%s\"", s)
	}
}
