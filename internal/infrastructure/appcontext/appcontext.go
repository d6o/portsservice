package appcontext

import (
	"context"
	"errors"

	"github.com/d6o/portsservice/internal/infrastructure/log"
	"go.uber.org/zap"
)

type indexContext int

const (
	loggerKey indexContext = iota
)

// WithLogger receives a loggerFromContext and stores it as the base loggerFromContext for the package.
func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// Logger returns a loggerFromContext with all values from context loaded on it.
func Logger(ctx context.Context) *zap.Logger {
	logger, err := loggerFromContext(ctx)
	if err != nil {
		logger = log.NewZapNop()
	}

	return logger
}

func loggerFromContext(ctx context.Context) (*zap.Logger, error) {
	logger, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok {
		return nil, errors.New("there is no loggerFromContext on context")
	}

	return logger, nil
}
