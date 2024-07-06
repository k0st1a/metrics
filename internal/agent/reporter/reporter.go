// Package reporter is send metrics to HTTP server.
package reporter

import (
	"context"
	"sync"
	"time"

	"github.com/k0st1a/metrics/internal/agent/report/json"
	"github.com/k0st1a/metrics/internal/pkg/rawmetric"
	"github.com/k0st1a/metrics/internal/ports"
	"github.com/rs/zerolog/log"
)

type state struct {
	client         ports.DoBatcher
	pollerCh       chan<- struct{}
	serverAddr     string
	reportInterval int
	rateLimit      int
}

// NewReporter - создание репортера, который отправляет метрики на сервер, где:
//   - serverAddr - адрес сервера;
//   - reportInterval - интервал между отправками на сервер, в секундах;
//   - rateLimit - количество одновременных запросов на сервер;
//   - sign - функция подписи передаваемых на сервер данных.
//
//nolint:lll //no need here
func NewReporter(serverAddr string, reportInterval int, rateLimit int, client ports.DoBatcher) (*state, <-chan struct{}) {
	pollerCh := make(chan struct{})
	return &state{
		serverAddr:     serverAddr,
		reportInterval: reportInterval,
		rateLimit:      rateLimit,
		pollerCh:       pollerCh,
		client:         client,
	}, pollerCh
}

// Do - запуск репортера, где:
//   - ctx - контекст отмены репортера;
//   - reportCh - канал получения метрик.
func (s *state) Do(ctx context.Context, reportCh <-chan map[string]rawmetric.MetricInfoRaw) {
	var wg sync.WaitGroup
	agentCh := make(chan map[string]rawmetric.MetricInfoRaw)
	for i := 0; i < s.rateLimit; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			json.NewReport(s.serverAddr, s.client, agentCh).Do(ctx)
		}()
	}

	reportTicker := time.NewTicker(time.Duration(s.reportInterval) * time.Second)

	for {
		select {
		case <-reportTicker.C:
			log.Printf("-->reportTick\n")
			s.pollerCh <- struct{}{}
			m := <-reportCh
			if len(m) == 0 {
				continue
			}
			agentCh <- m
		case <-ctx.Done():
			log.Printf("Reporter closed with cause:%s\n", ctx.Err())
			reportTicker.Stop()
			wg.Wait()
			return
		}
	}
}
