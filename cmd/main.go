package main

import (
	"context"
	"fmt"

	"github.com/palindrom615/requestbin"
	"github.com/palindrom615/requestbin/handler"
	"github.com/palindrom615/requestbin/handler/db"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const requestCtxKey = "request"

func HandleRequest(ctx context.Context, request *events.LambdaFunctionURLRequest) (res *events.LambdaFunctionURLResponse, e error) {
	logger := requestbin.GetLogger()
	var awsCfg, err = config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logger.Fatal(err)
		panic(err)
	}

	h := handler.NewConsHandler(
		handler.NewConsHandler(
			handler.NewEmbedCtxHandler(
				func(ctx context.Context, input *events.LambdaFunctionURLRequest) (handler.CtxKey, *events.LambdaFunctionURLRequest) {
					return requestCtxKey, input
				},
			),
			handler.NewFilteringHandler(
				func(ctx context.Context, req *events.LambdaFunctionURLRequest) bool {
					return req.RequestContext.HTTP.Method == "POST"
				},
			),
		),
		handler.NewConsHandler(
			handler.NewMappingHandler(
				func(ctx context.Context, input *events.LambdaFunctionURLRequest) (string, error) {

					return "", nil
				},
			),
			db.NewDynamoPutHandler(
				dynamodb.NewFromConfig(awsCfg),
				"host",
				func(ctx context.Context, input string) map[string]types.AttributeValue {
					m := make(map[string]types.AttributeValue)
					var e error
					m["mid"], _ = attributevalue.Marshal("key")
					m["info"], e = attributevalue.Marshal(input)
					if e != nil {
						logger.Errorw("info marshall fail", "err", e)
					}
					return m
				},
			),
		),
	)
	handlerCtx := context.Background()

	inputChan := make(chan *events.LambdaFunctionURLRequest)
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

	inputChan <- request
	<-done
	logger.Infow("done")

	return
}

func main() {
	lambda.Start(HandleRequest)
}
