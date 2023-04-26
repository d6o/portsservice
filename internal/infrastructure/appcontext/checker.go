package appcontext

import (
	"context"

	"go.uber.org/zap"
)

type (
	Checker struct{}
)

func NewChecker() *Checker {
	return &Checker{}
}

func (Checker) CheckContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		err := ctx.Err()
		Logger(ctx).With(zap.Error(err)).Warn("Context is closed, returning error")
		return err
	default:
		return nil
	}
}
