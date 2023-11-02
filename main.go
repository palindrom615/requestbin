package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var tableName = "hosts"

func HandleRequest(ctx context.Context, request events.LambdaFunctionURLRequest) *events.LambdaFunctionURLResponse {
	logger, _ := NewLogger()
	defer logger.Sync()
	logger.Infow("request", "request", request)
	db := NewDB(tableName, logger)
	err := db.PutValue(request.RequestContext.RequestID, request.Body)
	if err != nil {
		return &events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Headers:    map[string]string{"Content-Type": "text/plain"},
			Body:       "Internal Server Error",
		}
	}
	if request.RequestContext.HTTP.Method != "POST" {
		return &events.LambdaFunctionURLResponse{
			StatusCode: 405,
			Headers:    map[string]string{"Content-Type": "text/plain"},
			Body:       "Method Not Allowed",
		}
	}
	return &events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "text/plain"},
	}
}

func main() {
	lambda.Start(HandleRequest)
}
