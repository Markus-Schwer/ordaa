package config

import "github.com/caarlos0/env/v11"

type DatabaseConfig struct {
	URL string `env:"URL"`
}

func LoadDatabaseConfig() (*DatabaseConfig, error) {
	var cfg DatabaseConfig
	if err := env.ParseWithOptions(&cfg, env.Options{
		Prefix: "DATABASE_",
	}); err != nil {
		return nil, err
	}

	return &cfg, nil
}
