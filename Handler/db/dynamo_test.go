package db_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/palindrom615/requestbin/handler/db"
)

func TestNewDynamoPutHandler(t *testing.T) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
	)
	if err != nil {
		panic(err.Error())
	}

	dynamoClient := dynamodb.NewFromConfig(cfg)
	keyVal := func(ctx context.Context, input interface{}) map[string]types.AttributeValue {
		m := make(map[string]types.AttributeValue)
		m["mid"], _ = attributevalue.Marshal("key")
		m["info"], _ = attributevalue.Marshal("{}")
		return m
	}

	handler := db.NewDynamoPutHandler(
		dynamoClient,
		"host",
		keyVal,
	)
	inputChan := make(chan interface{})
	go func() {
		inputChan <- nil
	}()
	_, _, e := handler.Handle(context.Background(), inputChan)
	if e != nil {
		t.Error(e)
	}
}
