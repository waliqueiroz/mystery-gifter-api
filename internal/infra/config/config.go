package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Database DatabaseConfig
	Auth     AuthConfig
}

type DatabaseConfig struct {
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	Database string `env:"DB_DATABASE"`
	Username string `env:"DB_USERNAME"`
	Password string `env:"DB_PASSWORD"`
}

type AuthConfig struct {
	SecretKey       string        `env:"AUTH_SECRET_KEY"`
	SessionDuration time.Duration `env:"AUTH_SESSION_DURATION"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := env.ParseWithOptions(&cfg, env.Options{RequiredIfNoDef: true}); err != nil {
		return nil, err
	}

	return &cfg, nil
}
