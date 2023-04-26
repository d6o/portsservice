package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZap() (*zap.Logger, error) {
	zapConfig := zap.NewProductionConfig()
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := zapConfig.Build()
	return logger, err
}

func NewZapNop() *zap.Logger {
	return zap.NewNop()
}
