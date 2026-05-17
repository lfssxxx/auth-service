package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env      string        `envconfig:"ENV" default:"local"`
	TokenTTL time.Duration `envconfig:"TOKEN_TTL" required:"true"`
	Port     int           `envconfig:"PORT"`
	Timeout  time.Duration `envconfig:"TIMEOUT"`
}

func MustLoad() (Config, error) {

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, fmt.Errorf("config load: %w", err)
	}
	return cfg, nil
}
