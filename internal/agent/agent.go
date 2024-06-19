// Package agent - пакет HTTP-клиента для сбора рантайм-метрик и
// их последующей отправки на сервер по протоколу HTTP.
package agent

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/k0st1a/metrics/internal/agent/poller"
	"github.com/k0st1a/metrics/internal/agent/reporter"
	"github.com/k0st1a/metrics/internal/metrics/gopsutil"
	"github.com/k0st1a/metrics/internal/metrics/runtime"
	"github.com/k0st1a/metrics/internal/middleware/encrypt"
	"github.com/k0st1a/metrics/internal/middleware/roundtrip"
	"github.com/k0st1a/metrics/internal/middleware/sign"
	"github.com/k0st1a/metrics/internal/pkg/crypto/rsa"
	"github.com/k0st1a/metrics/internal/pkg/hash"
)

// Run - запуск агента.
func Run() error {
	cfg, err := NewConfig()
	if err != nil {
		return err
	}

	printConfig(cfg)

	ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelCtx()

	rm := runtime.NewMetric()
	gm := gopsutil.NewMetric()
	p, pc := poller.NewPoller(cfg.PollInterval, rm, gm)

	var middlewares []roundtrip.Middleware

	if cfg.HashKey != "" {
		h := hash.New(cfg.HashKey)
		middlewares = append(middlewares, sign.New(h))
	}

	if cfg.CryptoKey != "" {
		pbl, err := rsa.NewPublicFromFile(cfg.CryptoKey)
		if err != nil {
			return fmt.Errorf("rsa new public from file error:%w", err)
		}

		middlewares = append(middlewares, encrypt.New(pbl))
	}

	rt := roundtrip.New(http.DefaultTransport, middlewares...)

	r, rc := reporter.NewReporter(cfg.ServerAddr, cfg.ReportInterval, cfg.RateLimit, rt)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		p.Do(ctx, rc)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		r.Do(ctx, pc)
	}()

	<-ctx.Done()
	wg.Wait()

	return nil
}
