package agent

import (
	"github.com/k0st1a/metrics/internal/agent/report"
	"github.com/k0st1a/metrics/internal/metrics"
	"net/http"
)

func Run() error {
	var myMetrics = &metrics.MyStats{}
	var myClient = &http.Client{}

	cfg, err := collectConfig()
	if err != nil {
		return err
	}

	printConfig(cfg)

	go metrics.RunUpdateMetrics(myMetrics, cfg.PollInterval)
	report.RunReportMetrics(cfg.ServerAddr, myClient, myMetrics, cfg.ReportInterval)

	return nil
}
