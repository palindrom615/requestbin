package datastore

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/palindrom615/requestbin"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var logger = requestbin.GetLogger()

// Table is a struct that represents a DynamoDB datastore.
// name is the name of the datastore.
// keyName is the primary keyName of the datastore.
type Table struct {
	name    string
	keyName string
}

func NewTable(name string, keyName string) *Table {
	return &Table{
		name:    name,
		keyName: keyName,
	}
}

func (t *Table) MakePutItemInput(key string, values map[string]interface{}) (*dynamodb.PutItemInput, error) {
	items, err := attributevalue.MarshalMap(values)
	if err != nil {
		logger.Errorw("PutItem marshal failed", "error", err, "values", values)
		return nil, err
	}

	items[t.keyName] = &types.AttributeValueMemberS{Value: key}

	return &dynamodb.PutItemInput{
		TableName: aws.String(t.name),
		Item:      items,
	}, nil
}

func (t *Table) MakeUpdateItemInput(key string, values map[string]interface{}) (*dynamodb.UpdateItemInput, error) {
	keyMap, _ := attributevalue.MarshalMap(map[string]interface{}{
		t.keyName: key,
	})
	b := expression.NewBuilder()
	for k, v := range values {
		b = b.WithUpdate(expression.Set(expression.Name(k), expression.Value(v)))
	}
	e, err := b.Build()

	if err != nil {
		logger.Errorw("UpdateItem marshal failed", "error", err, "values", values)
		return nil, err
	}

	return &dynamodb.UpdateItemInput{
		TableName:                 aws.String(t.name),
		Key:                       keyMap,
		ExpressionAttributeNames:  e.Names(),
		ExpressionAttributeValues: e.Values(),
		UpdateExpression:          e.Update(),
	}, nil
}
