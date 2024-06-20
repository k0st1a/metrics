// Package reporter is send metrics to HTTP server.
package reporter

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/k0st1a/metrics/internal/agent/model"
	"github.com/k0st1a/metrics/internal/agent/report/json"
	"github.com/rs/zerolog/log"
)

type state struct {
	sign           http.RoundTripper
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
func NewReporter(serverAddr string, reportInterval int, rateLimit int, sign http.RoundTripper) (*state, <-chan struct{}) {
	pollerCh := make(chan struct{})
	return &state{
		serverAddr:     serverAddr,
		reportInterval: reportInterval,
		rateLimit:      rateLimit,
		pollerCh:       pollerCh,
		sign:           sign,
	}, pollerCh
}

// Do - запуск репортера, где:
//   - ctx - контекст отмены репортера;
//   - reportCh - канал получения метрик.
func (s *state) Do(ctx context.Context, reportCh <-chan map[string]model.MetricInfoRaw) {
	var wg sync.WaitGroup
	agentCh := make(chan map[string]model.MetricInfoRaw)
	for i := 0; i < s.rateLimit; i++ {
		c := &http.Client{
			Transport: s.sign,
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			json.NewReport(s.serverAddr, c, agentCh).Do(ctx)
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
