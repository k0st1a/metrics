// Package agent is for forces the client to send metrics
package agent

import (
	"context"

	"github.com/k0st1a/metrics/internal/pkg/rawmetric"
	"github.com/k0st1a/metrics/internal/ports"
	"github.com/rs/zerolog/log"
)

type state struct {
	client  ports.DoBatcher
	channel <-chan map[string]rawmetric.Info
}

// New - создание агента, который получает метрики через канал и заставляет клиента
// отправлять метрики на сервер, где:
//   - с - клиент;
//   - ch - через данный канал получаем метрики для отправки клиентом на сервер.
func New(c ports.DoBatcher, ch <-chan map[string]rawmetric.Info) *state {
	return &state{
		client:  c,
		channel: ch,
	}
}

// Do - запуск репортера.
func (s *state) Do(ctx context.Context) {
	for {
		select {
		case m := <-s.channel:
			s.client.DoBatch(rawmetric.Map2List(m))
		case <-ctx.Done():
			log.Printf("Reporter closed with cause:%s\n", ctx.Err())
			return
		}
	}
}
