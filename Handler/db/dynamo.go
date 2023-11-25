package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/palindrom615/requestbin"
)

type DynamoPutHandler struct {
	dynamoDbClient DynamoDBPutItemAPI
	tableName      string
	getItem        func(context context.Context, input interface{}) map[string]types.AttributeValue
}

type DynamoDBPutItemAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

func NewDynamoPutHandler(
	db DynamoDBPutItemAPI,
	tableName string,
	getItem func(context context.Context, input interface{}) map[string]types.AttributeValue,
) *DynamoPutHandler {
	return &DynamoPutHandler{
		dynamoDbClient: db,
		tableName:      tableName,
		getItem:        getItem,
	}
}

func (h *DynamoPutHandler) Handle(ctx context.Context, input <-chan interface{}) (context.Context, <-chan interface{}, error) {
	logger := requestbin.GetLogger()

	select {
	case i := <-input:
		p := h.getItem(ctx, i)
		logger.Infow("putItem: ", p)

		putItemInput := &dynamodb.PutItemInput{
			TableName: &h.tableName,
			Item:      p,
		}
		e := make(chan error)
		go func() {
			_, err := h.dynamoDbClient.PutItem(context.Background(), putItemInput)
			if err != nil {
				logger.Error(err)
			}
			e <- err
		}()
		return ctx, input, <-e
	case <-ctx.Done():
		return ctx, input, nil
	}

}
