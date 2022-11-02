package l

import (
	"go.uber.org/zap"
)

var gLogger *zap.Logger

func initLogger(prefix string) error {
	flags := new(LoggerAttr).InitDefaultLogger()
	flags.SetLogPath(prefix)
	logger, err := flags.NewLogger()
	if err != nil {
		logger.Error("init logger failed", zap.Error(err))
	}
	gLogger = logger
	return err
}

func GetLogger() *zap.Logger {
	return gLogger
}
