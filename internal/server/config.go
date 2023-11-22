package server

import (
	"errors"
	"flag"

	"github.com/caarlos0/env/v10"
	"github.com/k0st1a/metrics/internal/utils"
	"github.com/rs/zerolog/log"
)

type Config struct {
	ServerAddr string `env:"ADDRESS"`
}

func NewConfig() Config {
	return Config{
		ServerAddr: "localhost:8080",
	}
}

func parseEnv(cfg *Config) error {
	err := env.Parse(cfg)
	if err != nil {
		return err
	}
	return nil
}

func parseFlags(cfg *Config) error {
	addr := &utils.NetAddress{}
	addr.Set(cfg.ServerAddr)

	// если интерфейс не реализован,
	// здесь будет ошибка компиляции
	_ = flag.Value(addr)
	flag.Var(addr, "a", "server network address")
	flag.Parse()

	if len(flag.Args()) != 0 {
		return errors.New("unknown args")
	}

	cfg.ServerAddr = addr.String()

	return nil
}

func collectConfig(cfg *Config) error {
	log.Debug().
		Str("cfg.ServerAddr", cfg.ServerAddr).
		Msg("")

	err := parseFlags(cfg)
	if err != nil {
		return err
	}
	log.Debug().
		Str("cfg.ServerAddr", cfg.ServerAddr).
		Msg("After parseFlags")

	err = parseEnv(cfg)
	if err != nil {
		return err
	}
	log.Debug().
		Str("cfg.ServerAddr", cfg.ServerAddr).
		Msg("After parseEnv")

	return nil
}
