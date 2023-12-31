package requestbin

import "go.uber.org/zap"

var sugar *zap.SugaredLogger

func GetLogger() *zap.SugaredLogger {
	if sugar != nil {
		return sugar
	}
	logger, err := zap.NewDevelopment(zap.AddCaller())
	if err != nil {
		panic(err)
	}
	sugar = logger.Sugar()
	return sugar
}
