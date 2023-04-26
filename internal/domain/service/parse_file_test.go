package service_test

import (
	"context"
	"errors"
	"github.com/d6o/portsservice/internal/domain/model"
	"os"
	"testing"

	"github.com/d6o/portsservice/internal/domain/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestParsePorts_Parse(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := model.NewMockStorage(ctrl)
	storage.EXPECT().Save(ctx, "port_id_1", gomock.Any()).Return(nil)

	contextChecker := service.NewMockcontextChecker(ctrl)
	contextChecker.EXPECT().CheckContext(gomock.Any()).Return(nil).AnyTimes()

	parsePorts := service.NewParsePorts(storage, contextChecker)

	fileContent := `{
		"port_id_1": {
			"id": "port_id_1",
			"name": "Port 1",
			"city": "City 1",
			"country": "Country 1",
			"alias": ["Alias 1"],
			"regions": ["Region 1"],
			"coordinates": [1.0, 1.0],
			"province": "Province 1",
			"timezone": "Timezone 1",
			"unlocs": ["Unloc 1"],
			"code": "Code 1"
		}
	}`

	file, err := os.CreateTemp("", "ports.json")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	_, err = file.WriteString(fileContent)
	assert.NoError(t, err)

	err = file.Close()
	assert.NoError(t, err)

	err = parsePorts.Parse(context.Background(), file.Name())
	assert.NoError(t, err)
}

func TestParsePorts_Parse_StorageError(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := model.NewMockStorage(ctrl)
	storage.EXPECT().Save(ctx, "port_id_1", gomock.Any()).Return(errors.New("storage error"))

	contextChecker := service.NewMockcontextChecker(ctrl)
	contextChecker.EXPECT().CheckContext(gomock.Any()).Return(nil).AnyTimes()

	parsePorts := service.NewParsePorts(storage, contextChecker)

	fileContent := `{
		"port_id_1": {
			"id": "port_id_1",
			"name": "Port 1",
			"city": "City 1",
			"country": "Country 1",
			"alias": ["Alias 1"],
			"regions": ["Region 1"],
			"coordinates": [1.0, 1.0],
			"province": "Province 1",
			"timezone": "Timezone 1",
			"unlocs": ["Unloc 1"],
			"code": "Code 1"
		}
	}`

	file, err := os.CreateTemp("", "ports.json")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	_, err = file.WriteString(fileContent)
	assert.NoError(t, err)

	err = file.Close()
	assert.NoError(t, err)

	err = parsePorts.Parse(context.Background(), file.Name())
	assert.Error(t, err)
}

func TestParsePorts_ParseInvalidFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := model.NewMockStorage(ctrl)
	contextChecker := service.NewMockcontextChecker(ctrl)
	contextChecker.EXPECT().CheckContext(gomock.Any()).Return(nil).AnyTimes()

	parsePorts := service.NewParsePorts(storage, contextChecker)

	err := parsePorts.Parse(context.Background(), "non_existent_file.json")
	assert.Error(t, err)
}

func TestParsePorts_ParseInvalidJson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := model.NewMockStorage(ctrl)
	contextChecker := service.NewMockcontextChecker(ctrl)
	contextChecker.EXPECT().CheckContext(gomock.Any()).Return(nil).AnyTimes()

	parsePorts := service.NewParsePorts(storage, contextChecker)

	fileContent := `{
		"port_id_1": "Invalid JSON"
	}`

	file, err := os.CreateTemp("", "ports.json")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	_, err = file.WriteString(fileContent)
	assert.NoError(t, err)

	err = file.Close()
	assert.NoError(t, err)

	err = parsePorts.Parse(context.Background(), file.Name())
	assert.Error(t, err)
}

func TestParsePorts_CheckContextError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := model.NewMockStorage(ctrl)
	contextChecker := service.NewMockcontextChecker(ctrl)
	contextChecker.EXPECT().CheckContext(gomock.Any()).Return(errors.New("context error"))

	parsePorts := service.NewParsePorts(storage, contextChecker)

	err := parsePorts.Parse(context.Background(), "any_file.json")
	assert.Error(t, err)
	assert.Equal(t, errors.New("context error"), err)
}
