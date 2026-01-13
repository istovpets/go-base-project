package config

import (
	"errors"
	"log/slog"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	LogLevel  slog.LevelVar `env:"LOG_LEVEL" envDefault:"info"`
	Port      int16         `env:"PORT" envDefault:"9999"`
	DBConnStr string        `env:"DB_CONN_STR"`
}

func New() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.DBConnStr == "" {
		return errors.New("DB_CONN_STR must be set")
	}

	return nil
}
