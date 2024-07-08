// Package agent - пакет для сбора рантайм-метрик их последующей отправки на сервер.
package agent

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"

	grpcclient "github.com/k0st1a/metrics/internal/adapters/api/grpc/client"
	httpclient "github.com/k0st1a/metrics/internal/adapters/api/http/client"
	"github.com/k0st1a/metrics/internal/application/agent/config"
	"github.com/k0st1a/metrics/internal/pkg/agent"
	"github.com/k0st1a/metrics/internal/pkg/logwrap"
	"github.com/k0st1a/metrics/internal/pkg/metric/gopsutil"
	"github.com/k0st1a/metrics/internal/pkg/metric/runtime"
	"github.com/k0st1a/metrics/internal/pkg/poller"
	"github.com/k0st1a/metrics/internal/pkg/ratelimit"
	"github.com/k0st1a/metrics/internal/pkg/reporter"
	"github.com/k0st1a/metrics/internal/ports"
	"github.com/rs/zerolog/log"
)

// Run - запуск агента.
func Run() error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("make config error:%w", err)
	}

	err = logwrap.New(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("logwrap create error:%w", err)
	}

	log.Printf("Cfg:%+v", cfg)

	ctx, cancelFunc := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancelFunc()

	var c ports.DoBatcher

	switch cfg.APIType {
	case "http":
		c, err = httpclient.New(cfg)
		if err != nil {
			return fmt.Errorf("make client error:%w", err)
		}
	case "grpc":
		c, err = grpcclient.New(cfg)
		if err != nil {
			return fmt.Errorf("make client error:%w", err)
		}
	default:
		return fmt.Errorf("bad api type:%v", cfg.APIType)
	}

	rl := ratelimit.New(cfg.RateLimit)

	r, rc, clientCh := reporter.New(cfg.ReportInterval)

	rm := runtime.NewMetric()
	gm := gopsutil.NewMetric()
	p, pc := poller.NewPoller(cfg.PollInterval, rm, gm)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		rl.Run(func() {
			agent.New(c).Do(ctx, clientCh)
		})
	}()

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
