// Package agent is for forces the client to send metrics
package agent

import (
	"context"

	"github.com/k0st1a/metrics/internal/pkg/rawmetric"
	"github.com/k0st1a/metrics/internal/ports"
	"github.com/rs/zerolog/log"
)

type state struct {
	client ports.DoBatcher
}

// New - создание агента, который получает метрики через канал и заставляет клиента
// отправлять метрики на сервер, где:
//   - с - клиент;
func New(c ports.DoBatcher) *state {
	return &state{
		client: c,
	}
}

// Do - запуск репортера, где:
//   - ctx - контекст;
//   - ch - через данный канал получаем метрики для отправки клиентом на сервер.
func (s *state) Do(ctx context.Context, ch <-chan map[string]rawmetric.Info) {
	for {
		select {
		case m := <-ch:
			s.client.DoBatch(rawmetric.Map2List(m))
		case <-ctx.Done():
			log.Printf("Reporter closed with cause:%s\n", ctx.Err())
			return
		}
	}
}
