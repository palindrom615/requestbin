package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/palindrom615/requestbin"
)

var tableName = "hosts"

func HandleRequest(ctx context.Context, request events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLResponse, error) {
	logger, _ := requestbin.NewLogger()
	defer logger.Sync()
	logger.Infow("request", "request", request)
	db := requestbin.NewDB(tableName, logger)
	err := db.PutValue(request.RequestContext.RequestID, request.Body)
	if err != nil {
		return &events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Headers:    map[string]string{"Content-Type": "text/plain"},
			Body:       "Internal Server Error",
		}, nil
	}
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
