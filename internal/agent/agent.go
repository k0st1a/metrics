package agent

import (
	"net/http"

	"github.com/k0st1a/metrics/internal/agent/poller"
	"github.com/k0st1a/metrics/internal/agent/report/json"
	"github.com/k0st1a/metrics/internal/agent/reporter"
	"github.com/k0st1a/metrics/internal/metrics"
	"github.com/k0st1a/metrics/internal/middleware"
	"github.com/k0st1a/metrics/internal/utils"
)

func Run() error {
	cfg, err := collectConfig()
	if err != nil {
		return err
	}

	printConfig(cfg)

	h := utils.NewHash(cfg.HashKey)
	sgn := middleware.NewSign(http.DefaultTransport, h)

	var c = &http.Client{
		Transport: sgn,
	}

	m := metrics.NewMetrics()
	r := json.NewReport(cfg.ServerAddr, c, m)

	go poller.NewPoller(m, cfg.PollInterval).Run()
	reporter.NewReporter(r, cfg.ReportInterval).Run()

	return nil
}
