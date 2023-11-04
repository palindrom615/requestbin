package handler_test

import (
	"context"
	"testing"

	"github.com/palindrom615/requestbin/handler"
)

func TestLoggingHandler(t *testing.T) {
	// arrange
	l := handler.NewLoggingHandler[string]()
	i := make(chan string)

	//act
	ctx, o, e := l.Handle(context.Background(), i)

	go func() {
		i <- "test"
	}()

	// assert
	if e != nil {
		t.Errorf("LoggerHandler return non-nil error, got %v", e)
	}
	if handler.GetLogger(ctx) == nil {
		t.Errorf("LoggerHandler should attach logger to context")
	}
	if <-o != "test" {
		t.Errorf("LoggerHandler should return input as output 'test', got %v", i)
	}
}
