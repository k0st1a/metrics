package agent

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v10"
	"github.com/k0st1a/metrics/internal/utils"
	"github.com/rs/zerolog/log"
)

type Config struct {
	ServerAddr     string `env:"ADDRESS"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
}

func newConfig() *Config {
	return &Config{
		ServerAddr:     "localhost:8080",
		PollInterval:   2,
		ReportInterval: 10,
	}
}

func parseEnv(cfg *Config) error {
	err := env.Parse(cfg)
	if err != nil {
		return fmt.Errorf("env parse error:%w", err)
	}
	return nil
}

func parseFlags(cfg *Config) error {
	addr := &utils.NetAddress{}
	err := addr.Set(cfg.ServerAddr)
	if err != nil {
		return err
	}

	// если интерфейс не реализован,
	// здесь будет ошибка компиляции
	_ = flag.Value(addr)
	flag.Var(addr, "a", "server network address in a form host:port")

	flag.IntVar(&(cfg.PollInterval), "p", 2, "metrics polling rate to the server")
	flag.IntVar(&(cfg.ReportInterval), "r", 10, "frequency of sending metrics to the server")

	flag.Parse()
	cfg.ServerAddr = addr.String()

	if len(flag.Args()) != 0 {
		return fmt.Errorf("unknown args")
	}

	return nil
}

func collectConfig() (cfg *Config, err error) {
	cfg = newConfig()

	err = parseFlags(cfg)
	if err != nil {
		return nil, err
	}

	err = parseEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func printConfig(cfg *Config) {
	log.Debug().
		Str("cfg.ServerAddr", cfg.ServerAddr).
		Int("cfg.PollInterval", cfg.PollInterval).
		Int("cfg.ReportInterval", cfg.ReportInterval).
		Msg("")
}
