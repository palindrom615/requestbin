package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/palindrom615/requestbin"
)

type DynamoPutHandler struct {
	dynamoDbClient DynamoDBAPI
	tableName      string
	mapFunc        func(context context.Context, input interface{}) map[string]types.AttributeValue
}

type DynamoDBAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

func NewDynamoPutHandler(
	db DynamoDBAPI,
	tableName string,
	getItem func(context context.Context, input interface{}) map[string]types.AttributeValue,
) *DynamoPutHandler {
	return &DynamoPutHandler{
		dynamoDbClient: db,
		tableName:      tableName,
		mapFunc:        getItem,
	}
}

func (h *DynamoPutHandler) Handle(ctx context.Context, input <-chan interface{}) (context.Context, <-chan interface{}) {
	logger := requestbin.GetLogger()
	newCtx, cancelFunc := context.WithCancelCause(ctx)

	select {
	case i := <-input:
		output := make(chan any)
		p := h.mapFunc(ctx, i)
		putItemInput := &dynamodb.PutItemInput{
			TableName: &h.tableName,
			Item:      p,
		}
		logger.Debugw("putItem created", "putItemInput", putItemInput)

		go func() {
			defer close(output)
			r, err := h.dynamoDbClient.PutItem(context.Background(), putItemInput)
			logger.Debugw("dynamoDb.PutItem returns", "return", r)
			if err != nil {
				logger.Error(err)
				cancelFunc(err)
			}
			output <- r
		}()
		return newCtx, output
	case <-ctx.Done():
		cancelFunc(ctx.Err())
		return newCtx, input
	}

}
