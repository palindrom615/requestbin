package main

import (
	"context"

	"github.com/palindrom615/requestbin/handler"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var tableName = "hosts"

func HandleRequest(ctx context.Context, request events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLResponse, error) {

	h := handler.NewComposeHandler(
		handler.NewEmbedInputHandler(
			func(ctx context.Context, input interface{}) handler.CtxKey {
				return "request"
			},
		),
		handler.NewLoggingHandler[interface{}](),
	)
	handlerCtx := context.Background()
	inputChan := make(chan interface{})
	h.Handle(handlerCtx, inputChan)

	go func() {
		inputChan <- &request
	}()

	if request.RequestContext.HTTP.Method != "POST" {
		return &events.LambdaFunctionURLResponse{
			StatusCode: 405,
			Headers:    map[string]string{"Content-Type": "text/plain"},
			Body:       "Method Not Allowed",
		}, nil
	}
	return &events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "text/plain"},
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
