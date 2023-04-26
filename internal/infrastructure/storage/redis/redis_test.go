package redis_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/d6o/portsservice/internal/domain/model"
	"github.com/d6o/portsservice/internal/infrastructure/storage/redis"
	"github.com/golang/mock/gomock"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedis_SaveAndGet(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := redis.NewMockclient(ctrl)
	contextChecker := redis.NewMockcontextChecker(ctrl)

	r := redis.NewRedis(contextChecker, client)

	port := &model.Port{
		Name:        "Port 1",
		City:        "City 1",
		Country:     "Country 1",
		Alias:       []string{"Alias 1"},
		Regions:     []string{"Region 1"},
		Coordinates: []float64{1.0, 1.0},
		Province:    "Province 1",
		Timezone:    "Timezone 1",
		Unlocs:      []string{"Unloc 1"},
		Code:        "Code 1",
	}

	jsonValue, err := json.Marshal(port)
	assert.NoError(t, err)

	contextChecker.EXPECT().CheckContext(ctx).Return(nil).Times(2)
	client.EXPECT().Set(ctx, "1", jsonValue, gomock.Any()).Return(goredis.NewStatusResult("", nil))
	client.EXPECT().Get(ctx, "1").Return(goredis.NewStringResult(string(jsonValue), nil))

	err = r.Save(ctx, "1", port)
	assert.NoError(t, err)

	result, ok, err := r.Get(ctx, "1")
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, port, result)
}

func TestRedis_GetNotFound(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := redis.NewMockclient(ctrl)
	contextChecker := redis.NewMockcontextChecker(ctrl)
	contextChecker.EXPECT().CheckContext(ctx).Return(nil)

	r := redis.NewRedis(contextChecker, client)

	client.EXPECT().Get(ctx, "not_found").Return(goredis.NewStringResult("", goredis.Nil))

	_, ok, err := r.Get(ctx, "not_found")
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestRedis_CheckContextError(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := redis.NewMockclient(ctrl)
	contextChecker := redis.NewMockcontextChecker(ctrl)
	contextChecker.EXPECT().CheckContext(ctx).Return(errors.New("context error")).Times(2)

	r := redis.NewRedis(contextChecker, client)

	port := &model.Port{
		Name: "Test",
	}

	err := r.Save(context.Background(), "1", port)
	assert.Error(t, err)
	assert.Equal(t, errors.New("context error"), err)

	_, _, err = r.Get(context.Background(), "1")
	assert.Error(t, err)
	assert.Equal(t, errors.New("context error"), err)
}
