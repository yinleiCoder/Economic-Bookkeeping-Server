package initialize

import "go.uber.org/zap"

func InitLogger() {
	//logger, _ := zap.NewProduction()
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}
