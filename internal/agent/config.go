package agent

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v10"
	"github.com/k0st1a/metrics/internal/pkg/netaddr"
	"github.com/rs/zerolog/log"
)

const (
	defaultPollInterval   = 2
	defaultReportInterval = 10
	defaultServerAddr     = "localhost:8080"
	defaultHashKey        = ""
	defaultRateLimit      = 1
)

type Config struct {
	ServerAddr     string `env:"ADDRESS"`
	HashKey        string `env:"KEY"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

func newConfig() *Config {
	return &Config{
		ServerAddr:     defaultServerAddr,
		PollInterval:   defaultPollInterval,
		ReportInterval: defaultReportInterval,
		HashKey:        defaultHashKey,
		RateLimit:      defaultRateLimit,
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
	addr := &netaddr.NetAddress{}
	err := addr.Set(cfg.ServerAddr)
	if err != nil {
		return fmt.Errorf("addr set error:%w", err)
	}

	// если интерфейс не реализован,
	// здесь будет ошибка компиляции
	_ = flag.Value(addr)
	flag.Var(addr, "a", "server network address in a form host:port")

	flag.IntVar(&(cfg.PollInterval), "p", cfg.PollInterval, "metrics polling rate to the server")
	flag.IntVar(&(cfg.ReportInterval), "r", cfg.ReportInterval, "frequency of sending metrics to the server")
	flag.StringVar(&(cfg.HashKey), "k", cfg.HashKey,
		"Hash key with which the request body will be encoded"+
			"HTTP Header HashSHA256 will be added to the HTTP request")
	flag.IntVar(&(cfg.RateLimit), "l", cfg.RateLimit, "number of simultaneously outgoing requests to the server")

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
		Str("cfg.HashKey", cfg.HashKey).
		Int("cfg.RateLimit", cfg.RateLimit).
		Msg("")
}
