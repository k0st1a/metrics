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

	m := metrics.NewMetrics()

	go poller.NewPoller(m, cfg.PollInterval).Run()

	r, mc := reporter.NewReporter(m, cfg.ReportInterval)

	h := utils.NewHash(cfg.HashKey)
	sgn := middleware.NewSign(http.DefaultTransport, h)

	for i := 0; i < cfg.RateLimit; i++ {
		c := &http.Client{
			Transport: sgn,
		}
		go json.NewReport(cfg.ServerAddr, c, mc).Do()
	}

	r.Run()

	return nil
}
