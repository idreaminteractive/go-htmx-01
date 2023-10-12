package config

import (
	"github.com/caarlos0/env/v9"
)

type EnvConfig struct {
	DatabaseFileName string `env:"DATABASE_FILENAME" envDefault:"/litefs/potato.db"`
	GoPort           string `env:"GO_PORT" envDefault:"8079"`
	DopplerConfig    string `env:"DOPPLER_CONFIG"`
	SecretSauce      string `env:"SECRET_SAUCE"`
}

func Parse() *EnvConfig {

	config := EnvConfig{}
	err := env.Parse(&config)
	if err != nil {
		panic("Could not parse env")
	}
	return &config
}
