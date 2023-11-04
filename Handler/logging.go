package handler

import (
	"context"

	"go.uber.org/zap"
)

type LoggerContextKey string

const LoggerKey LoggerContextKey = LoggerContextKey("logger")

type LoggingHandler[I any] struct {
	logger *zap.SugaredLogger
}

func NewLoggingHandler[I any]() *LoggingHandler[I] {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()
	defer sugar.Sync()
	return &LoggingHandler[I]{
		logger: sugar,
	}
}

func GetLogger(ctx context.Context) *zap.SugaredLogger {
	return ctx.Value(LoggerKey).(*zap.SugaredLogger)
}

func (h *LoggingHandler[I]) Handle(ctx context.Context, input <-chan I) (context.Context, <-chan I, error) {
	newCtx := context.WithValue(ctx, LoggerKey, h.logger)
	logger := GetLogger(newCtx)
	logger.Info("Logger attached to context")
	return newCtx, input, nil
}
