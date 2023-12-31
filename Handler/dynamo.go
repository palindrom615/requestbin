package handler

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoPutHandler[I any] struct {
	dynamoDbClient DynamoDBAPI
	getInput       func(context context.Context, input I) (*dynamodb.PutItemInput, error)
}

type DynamoDBAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

func NewDynamoPutHandler[I any](
	db DynamoDBAPI,
	getInput func(context context.Context, input I) (*dynamodb.PutItemInput, error),
) Handler[I, *dynamodb.PutItemOutput] {
	return &DynamoPutHandler[I]{
		dynamoDbClient: db,
		getInput:       getInput,
	}
}

func (h *DynamoPutHandler[I]) Handle(ctx context.Context, input <-chan I) (context.Context, <-chan *dynamodb.PutItemOutput) {
	newCtx, cancel := context.WithCancelCause(ctx)
	output := make(chan *dynamodb.PutItemOutput)

	select {
	case i := <-input:
		p, e := h.getInput(ctx, i)
		if e != nil {
			logger.Error("PutItemInput creation failed", e)
			cancel(e)
			return newCtx, output
		}
		logger.Debugw("putItem created", "putItemInput", p)

		go func() {
			defer close(output)
			r, err := h.dynamoDbClient.PutItem(context.Background(), p)
			logger.Debugw("dynamoDb.PutItem returns", "return", r)
			if err != nil {
				logger.Error(err)
				cancel(err)
			}
			output <- r
		}()
		return newCtx, output
	case <-ctx.Done():
		cancel(context.Cause(ctx))
		close(output)
		return newCtx, output
	}

}
