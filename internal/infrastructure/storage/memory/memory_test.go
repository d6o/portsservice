package memory_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"testing"

	"github.com/d6o/portsservice/internal/domain/model"
	"github.com/d6o/portsservice/internal/infrastructure/storage/memory"
	"github.com/stretchr/testify/assert"
)

func TestInMemory_SaveAndGet(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	contextChecker := memory.NewMockcontextChecker(ctrl)
	contextChecker.EXPECT().CheckContext(ctx).Return(nil).Times(2)

	db := memory.NewInMemory(contextChecker)

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

	err := db.Save(ctx, "1", port)
	assert.NoError(t, err)

	result, ok, err := db.Get(ctx, "1")
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, port, result)
}

func TestInMemory_GetNotFound(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	contextChecker := memory.NewMockcontextChecker(ctrl)
	contextChecker.EXPECT().CheckContext(ctx).Return(nil)

	db := memory.NewInMemory(contextChecker)

	_, ok, err := db.Get(ctx, "not_found")
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestInMemory_CheckContextError(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	contextChecker := memory.NewMockcontextChecker(ctrl)
	contextChecker.EXPECT().CheckContext(ctx).Return(errors.New("context error")).Times(2)

	db := memory.NewInMemory(contextChecker)

	port := &model.Port{
		Name: "Test",
	}

	err := db.Save(ctx, "1", port)
	assert.Error(t, err)
	assert.Equal(t, errors.New("context error"), err)

	_, _, err = db.Get(ctx, "1")
	assert.Error(t, err)
	assert.Equal(t, errors.New("context error"), err)
}
