package log

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var zapLogger *zap.Logger

func init() {
	var config zap.Config

	config = zap.NewProductionConfig()

	config.DisableStacktrace = true

	if gin.IsDebugging() {
		zapLogger, _ = config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	} else {
		zapLogger, _ = config.Build(zap.AddCallerSkip(1))
	}

}

func ZapLogger() *zap.Logger {
	return zapLogger
}
