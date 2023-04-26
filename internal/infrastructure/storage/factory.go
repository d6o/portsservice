package storage

import (
	"context"

	"github.com/d6o/portsservice/internal/domain/model"
	"github.com/d6o/portsservice/internal/infrastructure/appcontext"
	"github.com/d6o/portsservice/internal/infrastructure/storage/memory"
	"github.com/d6o/portsservice/internal/infrastructure/storage/redis"
	"go.uber.org/zap"
)

type (
	contextChecker interface {
		CheckContext(ctx context.Context) error
	}
)

func New(
	ctx context.Context,
	contextChecker contextChecker,
	storageType string,
	redisAddress string,
) (model.Storage, error) {
	switch storageType {
	case "redis":
		redisClient, err := redis.Connect(ctx, redisAddress)
		if err != nil {
			return nil, err
		}
		return redis.NewRedis(contextChecker, redisClient), nil

	case "memory":
		return memory.NewInMemory(contextChecker), nil

	default:
		appcontext.Logger(ctx).With(zap.String("storage_type", storageType)).
			Warn("unsupported storage type - using in-memory storage")
		return memory.NewInMemory(contextChecker), nil
	}
}
