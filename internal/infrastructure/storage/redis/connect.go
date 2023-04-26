package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func Connect(ctx context.Context, addr string) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return redisClient, nil
}
