// Package json contains the reporter code which work with JSON format.
package json

import (
	"context"

	"github.com/k0st1a/metrics/internal/pkg/rawmetric"
	"github.com/k0st1a/metrics/internal/ports"
	"github.com/rs/zerolog/log"
)

type report struct {
	client  ports.DoBatcher
	channel <-chan map[string]rawmetric.MetricInfoRaw
	address string
}

// NewReport - создание репортера, HTTP клиента, отправляющего метрики в формате JSON, где:
//   - a - адрем сервера;
//   - с - HTTP клиент;
//   - ch - через данный канал получаем метрики для отправки на сервер.
func NewReport(a string, c ports.DoBatcher, ch <-chan map[string]rawmetric.MetricInfoRaw) *report {
	return &report{
		address: a,
		client:  c,
		channel: ch,
	}
}

// Do - запуск репортера.
func (r *report) Do(ctx context.Context) {
	for {
		select {
		case mi := <-r.channel:
			r.client.DoBatch(rawmetric.Map2List(mi))
		case <-ctx.Done():
			log.Printf("JSON peporter closed with cause:%s\n", ctx.Err())
			return
		}
	}
}
