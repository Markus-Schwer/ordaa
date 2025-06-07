package config

import "github.com/caarlos0/env/v11"

type LogConfig struct {
	Level string `env:"LEVEL" envDefault:"INFO"`
	JSON  bool   `env:"JSON"`
}

func LoadLogConfig() (*LogConfig, error) {
	var cfg LogConfig
	if err := env.ParseWithOptions(&cfg, env.Options{
		Prefix: "LOG_",
	}); err != nil {
		return nil, err
	}

	return &cfg, nil
}
