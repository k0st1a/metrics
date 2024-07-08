// Package reporter is send metrics to HTTP server.
package reporter

import (
	"context"
	"time"

	"github.com/k0st1a/metrics/internal/pkg/rawmetric"
	"github.com/rs/zerolog/log"
)

type state struct {
	pollerCh       chan<- struct{}
	clientCh       chan<- map[string]rawmetric.Info
	reportInterval int
}

// New - создание репортера, который отправляет метрики на сервер, где:
//   - reportInterval - интервал между отправками на сервер, в секундах;
//   - pollerCh - канал через который сообщаем poller-у, о том, чтобы он отправил
//     метрики в канал reportCh (см. функцид Do);
//   - clientCh - канал через который отправляем метрики на клиента(ов).
func New(reportInterval int) (*state, <-chan struct{}, <-chan map[string]rawmetric.Info) {
	pollerCh := make(chan struct{})
	clientCh := make(chan map[string]rawmetric.Info)

	return &state{
		reportInterval: reportInterval,
		pollerCh:       pollerCh,
		clientCh:       clientCh,
	}, pollerCh, clientCh
}

// Do - запуск репортера, где:
//   - ctx - контекст отмены репортера;
//   - reportCh - канал, через который получаем метрики.
func (s *state) Do(ctx context.Context, reportCh <-chan map[string]rawmetric.Info) {

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
			s.clientCh <- m
		case <-ctx.Done():
			log.Printf("Reporter closed with cause:%s\n", ctx.Err())
			reportTicker.Stop()
			return
		}
	}
}
