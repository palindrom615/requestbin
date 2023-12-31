package handler

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoUpdateHandler[I any] struct {
	dynamoDbClient *dynamodb.Client
	tableName      string
	mapFunc        func(context context.Context, input I) (key map[string]types.AttributeValue, attr map[string]types.AttributeValueUpdate)
}

func (d DynamoUpdateHandler[I]) Handle(ctx context.Context, input <-chan I) (context.Context, <-chan *dynamodb.UpdateItemOutput) {
	newCtx, cancel := context.WithCancelCause(ctx)
	output := make(chan *dynamodb.UpdateItemOutput)

	select {
	case i := <-input:
		k, attr := d.mapFunc(ctx, i)

		item := &dynamodb.UpdateItemInput{
			TableName:        &d.tableName,
			Key:              k,
			AttributeUpdates: attr,
		}
		go func() {
			res, err := d.dynamoDbClient.UpdateItem(ctx, item)
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

func NewDynamoUpdateHandler[I any](
	db *dynamodb.Client,
	tableName string,
	mapFunc func(context context.Context, input I) (key map[string]types.AttributeValue, attr map[string]types.AttributeValueUpdate),
) Handler[I, *dynamodb.UpdateItemOutput] {
	return &DynamoUpdateHandler[I]{
		dynamoDbClient: db,
		tableName:      tableName,
		mapFunc:        mapFunc,
	}
}
