package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	ExternalAPIToken string `env:"EXTERNAL_API_TOKEN,notEmpty,unset,file"`
	DBConnString     string `env:"DB_CONN_STRING,notEmpty,unset,file"`
	Secret           string `env:"SECRET,notEmpty,unset,file"`
}

func New() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
