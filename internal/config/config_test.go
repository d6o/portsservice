package config_test

import (
	"os"
	"testing"

	"github.com/d6o/portsservice/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLoad_DefaultValues(t *testing.T) {
	os.Clearenv()

	cfg, err := config.Load()
	assert.NoError(t, err)

	assert.Equal(t, "memory", cfg.StorageProvider)
	assert.Equal(t, 8080, cfg.APIPort)
	assert.Equal(t, "", cfg.RedisAddress)
}

func TestLoad_EnvironmentValues(t *testing.T) {
	assert.NoError(t, os.Setenv("STORAGE_PROVIDER", "redis"))
	assert.NoError(t, os.Setenv("API_PORT", "9000"))
	assert.NoError(t, os.Setenv("REDIS_ADDRESS", "localhost:6379"))

	cfg, err := config.Load()
	assert.NoError(t, err)

	assert.Equal(t, "redis", cfg.StorageProvider)
	assert.Equal(t, 9000, cfg.APIPort)
	assert.Equal(t, "localhost:6379", cfg.RedisAddress)

	os.Clearenv()
}

func TestLoad_InvalidEnvironmentValues(t *testing.T) {
	assert.NoError(t, os.Setenv("API_PORT", "invalid_port"))

	_, err := config.Load()
	assert.Error(t, err)

	os.Clearenv()
}
