package storage

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string `envconfig:"POSTGRES_HOST" required:"true"`
	Port     string `envconfig:"POSTGRES_PORT" default:"5432"`
	User     string `envconfig:"POSTGRES_USER" required:"true"`
	Password string `envconfig:"POSTGRES_PASSWORD" required:"true"`
	Database string `envconfig:"POSTGRES_DB" required:"true"`
}

func NewConfig() (Config, error) {

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, fmt.Errorf("config load: %w", err)
	}
	return cfg, nil
}

func NewConfigMust() Config {
	cfg, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get postgres connection pool config: %w", err)
		panic(err)
	}
	return cfg
}
