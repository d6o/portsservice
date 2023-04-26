package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/d6o/portsservice/internal/domain/model"
	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -source redis.go -destination mock_redis.go -package redis

type (
	Redis struct {
		client         client
		contextChecker contextChecker
	}

	contextChecker interface {
		CheckContext(ctx context.Context) error
	}

	client interface {
		Get(ctx context.Context, key string) *redis.StringCmd
		Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	}
)

func NewRedis(contextChecker contextChecker, client client) *Redis {
	return &Redis{
		client:         client,
		contextChecker: contextChecker,
	}
}

const (
	noExpiration = 0
)

func (r Redis) Save(ctx context.Context, key string, value *model.Port) error {
	if err := r.contextChecker.CheckContext(ctx); err != nil {
		return err
	}

	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, jsonValue, noExpiration).Err()
}

func (r Redis) Get(ctx context.Context, key string) (*model.Port, bool, error) {
	if err := r.contextChecker.CheckContext(ctx); err != nil {
		return nil, false, err
	}

	val, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	var port model.Port
	return &port, true, json.Unmarshal([]byte(val), &port)
}
