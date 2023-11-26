package handler_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/palindrom615/requestbin/handler"
)

func getRealDB() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
	)
	if err != nil {
		panic(err.Error())
	}

	return dynamodb.NewFromConfig(cfg)
}

type MockDB struct {
	putItem func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

func (m *MockDB) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return m.putItem(ctx, params, optFns...)
}

func getMockDBSuccess() handler.DynamoDBAPI {
	return &MockDB{
		putItem: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return &dynamodb.PutItemOutput{}, nil
		},
	}
}

func TestNewDynamoPutHandler_success_without_err(t *testing.T) {
	// arrange
	handler := handler.NewDynamoPutHandler(
		getMockDBSuccess(),
		"host",
		func(ctx context.Context, input interface{}) map[string]types.AttributeValue {
			m := make(map[string]types.AttributeValue)
			m["mid"], _ = attributevalue.Marshal("key")
			m["info"], _ = attributevalue.Marshal("{asdf}")
			return m
		},
	)

	// act
	ctx, o := handlerTest(handler, nil)
	t.Log(<-o)
	t.Log(ctx)
}
