package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/palindrom615/requestbin"
	"github.com/palindrom615/requestbin/handler"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const requestCtxKey = "request"

var ErrInvalidBody = errors.New("invalid body")
var ErrInvalidMethod = errors.New("invalid method")

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
				func(ctx context.Context, req *events.LambdaFunctionURLRequest) error {
					if req.RequestContext.HTTP.Method != "POST" {
						return ErrInvalidMethod
					}
					return nil
				},
			),
		),
		handler.NewConsHandler(
			handler.NewMappingHandler(
				func(ctx context.Context, input *events.LambdaFunctionURLRequest) (m map[string]any, err error) {
					body := input.Body
					m = make(map[string]any)
					err = json.Unmarshal([]byte(body), &m)
					if m["code"] == nil || m["mid"] == nil {
						logger.Errorw("no required fields", "unmarshalled", m, "body", body)
						return m, ErrInvalidBody
					}
					return
				},
			),
			handler.NewDynamoPutHandler(
				dynamodb.NewFromConfig(awsCfg),
				"host",
				func(ctx context.Context, input map[string]any) map[string]types.AttributeValue {
					m := make(map[string]types.AttributeValue)
					var e error
					m["mid"], _ = attributevalue.Marshal(input["mid"])
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
			logger.Errorw("canceled", "error", context.Cause(outCtx))
			res = &events.LambdaFunctionURLResponse{
				StatusCode: 500,
				Headers:    map[string]string{"Content-Type": "text/plain"},
				Body:       fmt.Sprint(context.Cause(outCtx)),
			}

		case out := <-outChan:
			logger.Debugw("output", "out", out)
			res = &events.LambdaFunctionURLResponse{
				StatusCode: 200,
				Headers:    map[string]string{"Content-Type": "text/plain"},
				Body:       fmt.Sprint(out),
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
