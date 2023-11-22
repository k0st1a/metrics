package agent

import (
	"errors"
	"flag"

	"github.com/caarlos0/env/v10"
	"github.com/k0st1a/metrics/internal/utils"
	"github.com/rs/zerolog/log"
)

type Config struct {
	ServerAddr     string `env:"ADDRESS"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
}

func NewConfig() Config {
	return Config{
		ServerAddr:     "localhost:8080",
		PollInterval:   2,
		ReportInterval: 10,
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
	flag.Var(addr, "a", "server network address in a form host:port")

	flag.IntVar(&(cfg.PollInterval), "p", 2, "metrics polling rate to the server")
	flag.IntVar(&(cfg.ReportInterval), "r", 10, "frequency of sending metrics to the server")

	flag.Parse()
	cfg.ServerAddr = addr.String()

	if len(flag.Args()) != 0 {
		return errors.New("unknown args")
	}

	return nil
}

func collectConfig(cfg *Config) error {
	log.Debug().
		Str("cfg.ServerAddr", cfg.ServerAddr).
		Int("cfg.PollInterval", cfg.PollInterval).
		Int("cfg.ReportInterval", cfg.ReportInterval).
		Msg("")

	err := parseFlags(cfg)
	if err != nil {
		return err
	}
	log.Debug().
		Str("cfg.ServerAddr", cfg.ServerAddr).
		Int("cfg.PollInterval", cfg.PollInterval).
		Int("cfg.ReportInterval", cfg.ReportInterval).
		Msg("After parseFlags")

	err = parseEnv(cfg)
	if err != nil {
		return err
	}
	log.Debug().
		Str("cfg.ServerAddr", cfg.ServerAddr).
		Int("cfg.PollInterval", cfg.PollInterval).
		Int("cfg.ReportInterval", cfg.ReportInterval).
		Msg("After parseEnv")

	return nil
}
