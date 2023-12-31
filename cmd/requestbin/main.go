package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/palindrom615/requestbin"
	"github.com/palindrom615/requestbin/datastore"
	"github.com/palindrom615/requestbin/handler"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const requestCtxKey = "request"

var ErrInvalidBody = errors.New("invalid body")
var ErrInvalidMethod = errors.New("invalid method")
var logger = requestbin.GetLogger()

const tableName = "host"
const keyName = "mid"

func provideTable(tableName, keyName string) *datastore.Table {
	return datastore.NewTable(tableName, keyName)
}

func provideDynamoDBClient() *dynamodb.Client {
	var awsCfg, err = config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(err)
	}
	return dynamodb.NewFromConfig(awsCfg)
}

func provideHandler(cli *dynamodb.Client, table *datastore.Table) handler.Handler[*events.LambdaFunctionURLRequest, *dynamodb.PutItemOutput] {
	filteringMethod := func(ctx context.Context, req *events.LambdaFunctionURLRequest) error {
		if req.RequestContext.HTTP.Method != "POST" {
			return ErrInvalidMethod
		}
		return nil
	}
	bodyMapper := func(ctx context.Context, input *events.LambdaFunctionURLRequest) (m map[string]any, err error) {
		body := input.Body
		m = make(map[string]any)
		err = json.Unmarshal([]byte(body), &m)
		if m["code"] == nil || m["mid"] == nil {
			logger.Errorw("no required fields", "unmarshalled", m, "body", body)
			return m, ErrInvalidBody
		}
		return
	}
	dynamoAttributeMapper := func(context context.Context, input map[string]any) (*dynamodb.PutItemInput, error) {
		return table.MakePutItemInput(input["mid"].(string), input)
	}

	return handler.NewConsHandler(
		handler.NewConsHandler(
			handler.NewEmbedCtxHandler(
				func(ctx context.Context, input *events.LambdaFunctionURLRequest) (handler.CtxKey, *events.LambdaFunctionURLRequest) {
					return requestCtxKey, input
				},
			),
			handler.NewFilteringHandler(filteringMethod),
		),
		handler.NewConsHandler(
			handler.NewMappingHandler(bodyMapper),
			handler.NewDynamoPutHandler(
				cli,
				dynamoAttributeMapper,
			),
		),
	)
}

func provideHandlerRequest[O any](h handler.Handler[*events.LambdaFunctionURLRequest, O]) func(ctx context.Context, request *events.LambdaFunctionURLRequest) (res *events.LambdaFunctionURLResponse, e error) {
	return func(ctx context.Context, request *events.LambdaFunctionURLRequest) (res *events.LambdaFunctionURLResponse, e error) {
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
}

func main() {
	cli := provideDynamoDBClient()
	table := provideTable(tableName, keyName)
	h := provideHandler(cli, table)
	handleRequest := provideHandlerRequest(h)
	lambda.Start(handleRequest)
}
