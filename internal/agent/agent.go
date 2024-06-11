// Пакет HTTP-клиента для сбора рантайм-метрик и их последующей отправки на сервер по протоколу HTTP.
package agent

import (
	"net/http"

	"github.com/k0st1a/metrics/internal/agent/poller"
	"github.com/k0st1a/metrics/internal/agent/reporter"
	"github.com/k0st1a/metrics/internal/metrics/gopsutil"
	"github.com/k0st1a/metrics/internal/metrics/runtime"
	"github.com/k0st1a/metrics/internal/middleware"
	"github.com/k0st1a/metrics/internal/pkg/hash"
)

// Run - запуск агента.
func Run() error {
	cfg, err := NewConfig()
	if err != nil {
		return err
	}

	printConfig(cfg)

	rm := runtime.NewMetric()
	gm := gopsutil.NewMetric()
	p, pc := poller.NewPoller(cfg.PollInterval, rm, gm)

	// sign
	h := hash.New(cfg.HashKey)
	sign := middleware.NewSign(http.DefaultTransport, h)

	r, rc := reporter.NewReporter(cfg.ServerAddr, cfg.ReportInterval, cfg.RateLimit, sign)

	go p.Do(rc)

	r.Do(pc)

	return nil
}
