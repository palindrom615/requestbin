package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/palindrom615/requestbin"
)

type DynamoPutHandler struct {
	dynamoDbClient    DynamoDBAPI
	tableName         string
	mapInputToPutItem func(context context.Context, input interface{}) map[string]types.AttributeValue
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
		dynamoDbClient:    db,
		tableName:         tableName,
		mapInputToPutItem: getItem,
	}
}

func (h *DynamoPutHandler) Handle(ctx context.Context, input <-chan interface{}) (context.Context, <-chan interface{}) {
	logger := requestbin.GetLogger()
	newCtx, cancelFunc := context.WithCancelCause(ctx)

	select {
	case i := <-input:
		p := h.mapInputToPutItem(ctx, i)
		putItemInput := &dynamodb.PutItemInput{
			TableName: &h.tableName,
			Item:      p,
		}
		logger.Debugw("putItem created", "putItemInput", putItemInput)

		go func() {
			_, err := h.dynamoDbClient.PutItem(context.Background(), putItemInput)
			if err != nil {
				logger.Error(err)
				cancelFunc(err)
			}
		}()
		return newCtx, input
	case <-ctx.Done():
		cancelFunc(ctx.Err())
		return newCtx, input
	}

}
