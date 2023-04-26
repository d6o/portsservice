package config

import (
	"github.com/kelseyhightower/envconfig"
)

type (
	Specification struct {
		StorageProvider string `default:"memory" split_words:"true"`
		APIPort         int    `default:"8080" split_words:"true"`
		RedisAddress    string `split_words:"true"`
	}
)

const (
	noPrefix = ""
)

func Load() (*Specification, error) {
	var config Specification
	return &config, envconfig.Process(noPrefix, &config)
}
