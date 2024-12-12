package config

import "github.com/caarlos0/env/v11"

type Config struct {
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	Database string `env:"DB_DATABASE"`
	Username string `env:"DB_USERNAME"`
	Password string `env:"DB_PASSWORD"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
