package agent

import (
	"fmt"
	"net/http"

	"github.com/k0st1a/metrics/internal/agent/reporter"
	"github.com/k0st1a/metrics/internal/agent/reporter/json"
	"github.com/k0st1a/metrics/internal/metrics"
)

func Run() error {
	var ms = &metrics.MyStats{}
	var c = &http.Client{}

	cfg, err := collectConfig()
	if err != nil {
		return err
	}

	printConfig(cfg)

	tr, err := json.NewReporter(cfg.ServerAddr)
	if err != nil {
		return fmt.Errorf("json.NewReporter error:%w", err)
	}
	r := reporter.NewReporter(tr)

	go metrics.RunUpdateMetrics(ms, cfg.PollInterval)
	r.RunReportMetrics(c, ms, cfg.ReportInterval)

	return nil
}
