package main

import "go.uber.org/zap"

func NewLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction()
	sugar := logger.Sugar()
	return sugar, err
}
