package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// initialize logger
func InitLogger() {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Use ISO8601 time format

	// Disable JSON encoding
	cfg.Encoding = "console"

	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	Logger = l
}
