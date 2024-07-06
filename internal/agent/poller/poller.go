// Package poller for poll metrics from system.
package poller

import (
	"context"
	"sync"
	"time"

	"github.com/k0st1a/metrics/internal/pkg/collector"
	"github.com/k0st1a/metrics/internal/pkg/rawmetric"
	"github.com/k0st1a/metrics/internal/ports"
	"github.com/rs/zerolog/log"
)

type state struct {
	runtimeMetrics  ports.RawMetricInfoer
	gopsutilMetrics ports.RawMetricInfoer
	reportCh        chan<- map[string]rawmetric.Info
	pollInterval    int
}

// NewPoller - создание поллера, опросника метрик, где:
//   - i - через заданное количество секунд запускать сбор метрик;
//   - rm - функция формирования runtime метрик;
//   - gm - функция формирования gopsutil метрик.
func NewPoller(i int, rm ports.RawMetricInfoer, gm ports.RawMetricInfoer) (*state, <-chan map[string]rawmetric.Info) {
	reportCh := make(chan map[string]rawmetric.Info)
	return &state{
		pollInterval:    i,
		runtimeMetrics:  rm,
		gopsutilMetrics: gm,
		reportCh:        reportCh,
	}, reportCh
}

// Do - запуск опросника метрик, где:
//   - ctx - контест отмены опросника;
//   - reporterCh - признак того, что нужно отправить данные в канал reportCH
func (s *state) Do(ctx context.Context, reporterCh <-chan struct{}) {
	pollTicker := time.NewTicker(time.Duration(s.pollInterval) * time.Second)
	var wg sync.WaitGroup

	// runtime
	collectRuntimeCh := make(chan struct{}, 1)
	rcl, pollRuntimeCh := collector.NewCollector(collectRuntimeCh, s.runtimeMetrics)
	wg.Add(1)
	go func() {
		defer wg.Done()
		rcl.Do(ctx)
	}()

	// gopsutil
	collectGopsutilCh := make(chan struct{}, 1)
	gcl, pollGopsutilCh := collector.NewCollector(collectGopsutilCh, s.gopsutilMetrics)
	wg.Add(1)
	go func() {
		defer wg.Done()
		gcl.Do(ctx)
	}()

	acc := make(map[string]rawmetric.Info)

	for {
		select {
		case <-pollTicker.C:
			log.Printf("-->pollTick\n")
			collectRuntimeCh <- struct{}{}
			collectGopsutilCh <- struct{}{}

			rm, ok := <-pollRuntimeCh
			log.Printf("ok:%v, rm:%v\n", ok, rm)
			if ok {
				acc = rawmetric.Append(acc, rm)
			}

			gm, ok := <-pollGopsutilCh
			log.Printf("ok:%v, gm:%v\n", ok, gm)
			if ok {
				acc = rawmetric.Append(acc, gm)
			}
			log.Printf("acc after poll:%v\n", acc)
		case <-reporterCh:
			log.Printf("<-reportCh, acc:%v\n", acc)
			s.reportCh <- acc
			acc = map[string]rawmetric.Info{}
		case <-ctx.Done():
			log.Printf("Poller closed with cause:%s\n", ctx.Err())
			pollTicker.Stop()
			wg.Wait()
			return
		}
	}
}
