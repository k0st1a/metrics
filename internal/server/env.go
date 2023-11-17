package server

import (
	"github.com/caarlos0/env/v10"
)

func parseEnv(cfg *Config) error {
	err := env.Parse(cfg)
	if err != nil {
		return err
	}
	return nil
}
