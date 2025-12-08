package config

import (
	"log/slog"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	LogLevel slog.LevelVar `env:"LOG_LEVEL" envDefault:"info"`
}

func New() (*Config, error) {
	cfg := &Config{}

	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
