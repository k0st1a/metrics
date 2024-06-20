// Package poller for poll metrics from system.
package poller

import (
	"context"
	"sync"
	"time"

	"github.com/k0st1a/metrics/internal/agent/collector"
	"github.com/k0st1a/metrics/internal/agent/model"
	"github.com/rs/zerolog/log"
)

// MetricInfoRawer - интерфейс формирования метрик.
type MetricInfoRawer interface {
	MetricInfoRaw() []model.MetricInfoRaw
}

type state struct {
	runtimeMetrics  MetricInfoRawer
	gopsutilMetrics MetricInfoRawer
	reportCh        chan<- map[string]model.MetricInfoRaw
	pollInterval    int
}

// NewPoller - создание поллера, опросника метрик, где:
//   - i - через заданное количество секунд запускать сбор метрик;
//   - rm - функция формирования runtime метрик;
//   - gm - функция формирования gopsutil метрик.
func NewPoller(i int, rm MetricInfoRawer, gm MetricInfoRawer) (*state, <-chan map[string]model.MetricInfoRaw) {
	reportCh := make(chan map[string]model.MetricInfoRaw)
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

	acc := make(map[string]model.MetricInfoRaw)

	for {
		select {
		case <-pollTicker.C:
			log.Printf("-->pollTick\n")
			collectRuntimeCh <- struct{}{}
			collectGopsutilCh <- struct{}{}

			rm, ok := <-pollRuntimeCh
			log.Printf("ok:%v, rm:%v\n", ok, rm)
			if ok {
				acc = model.Append(acc, rm)
			}

			gm, ok := <-pollGopsutilCh
			log.Printf("ok:%v, gm:%v\n", ok, gm)
			if ok {
				acc = model.Append(acc, gm)
			}
			log.Printf("acc after poll:%v\n", acc)
		case <-reporterCh:
			log.Printf("<-reportCh, acc:%v\n", acc)
			s.reportCh <- acc
			acc = map[string]model.MetricInfoRaw{}
		case <-ctx.Done():
			log.Printf("Poller closed with cause:%s\n", ctx.Err())
			pollTicker.Stop()
			wg.Wait()
			return
		}
	}
}
