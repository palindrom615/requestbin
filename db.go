package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"go.uber.org/zap"
)

type DB struct {
	tableName string
	logger    *zap.SugaredLogger
}

func NewDB(tableName string, logger *zap.SugaredLogger) *DB {
	return &DB{tableName, logger}
}

func (db *DB) PutValue(key string, val string) error {
	db.logger.Infow("put key", "key", key, "val", val)
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	putItemInput := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      map[string]*dynamodb.AttributeValue{},
	}
	_, err := svc.PutItem(putItemInput)
	if err != nil {
		db.logger.Errorw("put key failed", "err", err)
		return err
	}
	return nil
}
