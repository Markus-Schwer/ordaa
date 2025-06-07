package config

import "github.com/caarlos0/env/v11"

type MatrixConfig struct {
	HomeserverURL string   `env:"HOMESERVER"`
	Username      string   `env:"USERNAME"`
	Password      string   `env:"PASSWORD"`
	Rooms         []string `env:"ROOMS"`
	DisplayName   string   `env:"DISPLAY_NAME" envDefault:"Chicken Masalla legende Wollmilchsau [BOT]"`
}

func LoadMatrixConfig() (*MatrixConfig, error) {
	var cfg MatrixConfig
	if err := env.ParseWithOptions(&cfg, env.Options{
		Prefix: "MATRIX_",
	}); err != nil {
		return nil, err
	}

	return &cfg, nil
}
