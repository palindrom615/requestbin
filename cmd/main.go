package main

import (
	"context"

	"github.com/palindrom615/requestbin"
	"github.com/palindrom615/requestbin/handler"
	"github.com/palindrom615/requestbin/handler/db"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func HandleRequest(ctx context.Context, request events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLResponse, error) {
	logger := requestbin.GetLogger()
	var awsCfg, err = config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logger.Fatal(err)
		panic(err)
	}

	dynamoDbHandler := buildDynamoDbHandler(awsCfg, "host")
	h := handler.NewComposeHandler(
		dynamoDbHandler,
		handler.NewEmbedInputHandler(
			func(ctx context.Context, input interface{}) handler.CtxKey {
				return "request"
			},
		),
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

func buildDynamoDbHandler(awsCfg aws.Config, tableNme string) *db.DynamoPutHandler {
	dynamoClient := dynamodb.NewFromConfig(awsCfg)
	keyVal := func(ctx context.Context, input interface{}) map[string]types.AttributeValue {
		m := make(map[string]types.AttributeValue)
		m["mid"], _ = attributevalue.Marshal("key")
		m["info"], _ = attributevalue.Marshal("{}")
		return m
	}
	return db.NewDynamoPutHandler(
		dynamoClient,
		tableNme,
		keyVal,
	)
}

func main() {
	lambda.Start(HandleRequest)
}
