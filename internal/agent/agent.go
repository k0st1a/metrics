package agent

import (
	"net/http"
	"time"

	"github.com/k0st1a/metrics/internal/agent/collector"
	"github.com/k0st1a/metrics/internal/agent/model"
	"github.com/k0st1a/metrics/internal/agent/report/json"
	"github.com/k0st1a/metrics/internal/metrics/gopsutil"
	"github.com/k0st1a/metrics/internal/metrics/runtime"
	"github.com/k0st1a/metrics/internal/middleware"
	"github.com/k0st1a/metrics/internal/utils"
)

func Run() error {
	cfg, err := collectConfig()
	if err != nil {
		return err
	}

	printConfig(cfg)

	// runtime
	rm := runtime.NewMetric()
	pollRCh := make(chan struct{}, 1)
	rcl, pollResRCh := collector.NewCollector(pollRCh, rm)
	go rcl.Do()

	// gopsutil
	gm := gopsutil.NewMetric()
	pollGCh := make(chan struct{}, 1)
	gcl, pollResGCh := collector.NewCollector(pollGCh, gm)
	go gcl.Do()

	// sign
	h := utils.NewHash(cfg.HashKey)
	sgn := middleware.NewSign(http.DefaultTransport, h)

	// reporters
	reportCh := make(chan []model.MetricInfo)
	for i := 0; i < cfg.RateLimit; i++ {
		c := &http.Client{
			Transport: sgn,
		}
		go json.NewReport(cfg.ServerAddr, c, reportCh).Do()
	}

	acc := []model.MetricInfo{}

	pollTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)

	for {
		select {
		case <-pollTicker.C:
			pollRCh <- struct{}{}
			pollGCh <- struct{}{}

			acc = []model.MetricInfo{}

			rm, ok := <-pollResRCh
			if ok {
				acc = append(acc, rm...)
			}

			gm, ok := <-pollResGCh
			if ok {
				acc = append(acc, gm...)
			}
		case <-reportTicker.C:
			reportCh <- acc
		}
	}
}
