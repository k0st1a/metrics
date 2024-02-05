package agent

import (
	"net/http"

	"github.com/k0st1a/metrics/internal/agent/poller"
	"github.com/k0st1a/metrics/internal/agent/report/json"
	"github.com/k0st1a/metrics/internal/agent/reporter"
	"github.com/k0st1a/metrics/internal/metrics"
)

func Run() error {
	var c = &http.Client{}

	cfg, err := collectConfig()
	if err != nil {
		return err
	}

	printConfig(cfg)

	m := metrics.NewMetrics()
	r := json.NewReport(cfg.ServerAddr, c, m)

	go poller.NewPoller(m, cfg.PollInterval).Run()
	reporter.NewReporter(r, cfg.ReportInterval).Run()

	return nil
}
