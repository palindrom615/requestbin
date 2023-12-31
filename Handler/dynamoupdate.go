package handler

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoUpdateHandler[I any] struct {
	dynamoDbClient *dynamodb.Client
	getInput       func(context context.Context, input I) (*dynamodb.UpdateItemInput, error)
}

func NewDynamoUpdateHandler[I any](
	db *dynamodb.Client,
	getInput func(context context.Context, input I) (*dynamodb.UpdateItemInput, error),
) Handler[I, *dynamodb.UpdateItemOutput] {
	return &DynamoUpdateHandler[I]{
		dynamoDbClient: db,
		getInput:       getInput,
	}
}

func (d DynamoUpdateHandler[I]) Handle(ctx context.Context, input <-chan I) (context.Context, <-chan *dynamodb.UpdateItemOutput) {
	newCtx, cancel := context.WithCancelCause(ctx)
	output := make(chan *dynamodb.UpdateItemOutput)

	select {
	case i := <-input:
		u, e := d.getInput(ctx, i)
		if e != nil {
			logger.Error("PutItemInput creation failed", e)
			cancel(e)
			return newCtx, output
		}
		logger.Debugw("updateItem created", "updateItemInput", u)

		go func() {
			defer close(output)
			res, err := d.dynamoDbClient.UpdateItem(ctx, u)
			logger.Debugw("dynamoDb.UpdateItem returns", "return", res)
			if err != nil {
				logger.Error(err)
				cancel(err)
			}
			output <- res
		}()
		return newCtx, output
	case <-ctx.Done():
		cancel(context.Cause(ctx))
		close(output)
		return newCtx, output
	}
}
