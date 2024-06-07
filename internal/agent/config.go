package agent

import (
	"flag"
	"fmt"
	"os"
	"strconv"

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
	ServerAddr     string
	HashKey        string
	PollInterval   int
	ReportInterval int
	RateLimit      int
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	addr := &netaddr.NetAddress{}
	err := addr.Set(defaultServerAddr)
	if err != nil {
		return nil, fmt.Errorf("addr set error:%w", err)
	}

	flag.Var(addr, "a", "server network address")

	flag.IntVar(&cfg.PollInterval, "p", defaultPollInterval, "metrics polling rate to the server")
	flag.IntVar(&cfg.ReportInterval, "r", defaultReportInterval, "frequency of sending metrics to the server")
	flag.StringVar(&cfg.HashKey, "k", defaultHashKey,
		"Hash key with which the request body will be encoded"+
			"HTTP Header HashSHA256 will be added to the HTTP request")
	flag.IntVar(&(cfg.RateLimit), "l", defaultRateLimit, "number of simultaneously outgoing requests to the server")

	flag.Parse()

	if len(flag.Args()) != 0 {
		return nil, fmt.Errorf("unknown args:%v", flag.Args())
	}

	cfg.ServerAddr = addr.String()

	sa, ok := os.LookupEnv("ADDRESS")
	if ok {
		cfg.ServerAddr = sa
	}

	k, ok := os.LookupEnv("KEY")
	if ok {
		cfg.HashKey = k
	}

	pi, ok := os.LookupEnv("POLL_INTERVAL")
	if ok {
		piInt, err := strconv.Atoi(pi)
		if err != nil {
			return nil, fmt.Errorf("POLL_INTERVAL parse error:%w", err)
		}

		cfg.PollInterval = piInt
	}

	ri, ok := os.LookupEnv("REPORT_INTERVAL")
	if ok {
		riInt, err := strconv.Atoi(ri)
		if err != nil {
			return nil, fmt.Errorf("REPORT_INTERVAL parse error:%w", err)
		}

		cfg.ReportInterval = riInt
	}

	rl, ok := os.LookupEnv("RATE_LIMIT")
	if ok {
		rlInt, err := strconv.Atoi(rl)
		if err != nil {
			return nil, fmt.Errorf("RATE_LIMIT parse error:%w", err)
		}

		cfg.RateLimit = rlInt
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
