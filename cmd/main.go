package main

import (
	"context"
	"fmt"

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

const requestCtxKey = "request"

func HandleRequest(ctx context.Context, request events.LambdaFunctionURLRequest) (res *events.LambdaFunctionURLResponse, e error) {
	logger := requestbin.GetLogger()
	var awsCfg, err = config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logger.Fatal(err)
		panic(err)
	}

	h := handler.NewComposeHandler(
		handler.NewEmbedCtxHandler(
			func(ctx context.Context, input interface{}) (handler.CtxKey, any) {
				return requestCtxKey, input
			},
		),
		buildDynamoDbHandler(awsCfg, "host"),
	)
	handlerCtx := context.Background()

	inputChan := make(chan interface{})
	defer close(inputChan)

	done := make(chan any)
	defer close(done)

	go func() {
		outCtx, outChan := h.Handle(handlerCtx, inputChan)
		select {
		case <-outCtx.Done():
			logger.Errorw("canceled", "error", outCtx.Err())
			res = &events.LambdaFunctionURLResponse{
				StatusCode: 500,
				Headers:    map[string]string{"Content-Type": "text/plain"},
				Body:       outCtx.Err().Error(),
			}

		case out := <-outChan:
			logger.Debugw("output", "out", out)
			res = &events.LambdaFunctionURLResponse{
				StatusCode: 200,
				Headers:    map[string]string{"Content-Type": "text/plain"},
				Body:       fmt.Sprintf("%s", out),
			}
		}
		done <- nil
	}()

	inputChan <- &request
	<-done
	logger.Infow("done")

	return
}

func buildDynamoDbHandler(awsCfg aws.Config, tableNme string) *db.DynamoPutHandler {
	dynamoClient := dynamodb.NewFromConfig(awsCfg)
	logger := requestbin.GetLogger()
	mapFunc := func(ctx context.Context, input interface{}) map[string]types.AttributeValue {
		m := make(map[string]types.AttributeValue)
		var e error
		m["mid"], _ = attributevalue.Marshal("key")
		m["info"], e = attributevalue.Marshal(input)
		if e != nil {
			logger.Errorw("info marshall fail", "err", e)
		}
		return m
	}
	return db.NewDynamoPutHandler(
		dynamoClient,
		tableNme,
		mapFunc,
	)
}

func main() {
	lambda.Start(HandleRequest)
}
