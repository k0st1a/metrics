package agent

import (
	"github.com/k0st1a/metrics/internal/agent/report"
	"github.com/k0st1a/metrics/internal/metrics"
	"github.com/rs/zerolog/log"
	"net/http"
)

func Run() error {
	var myMetrics = &metrics.MyStats{}
	var myClient = &http.Client{}

	cfg := NewConfig()
	err := collectConfig(&cfg)
	if err != nil {
		return err
	}

	go metrics.RunUpdateMetrics(myMetrics, cfg.PollInterval)
	report.RunReportMetrics(cfg.ServerAddr, myClient, myMetrics, cfg.ReportInterval)

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
