// Package collector for collect metrics.
package collector

import (
	"context"

	"github.com/k0st1a/metrics/internal/agent/model"
	"github.com/rs/zerolog/log"
)

// MetricInfoRawer - интерфейс формирования метрик.
type MetricInfoRawer interface {
	MetricInfoRaw() []model.MetricInfoRaw
}

type state struct {
	in     <-chan struct{}
	out    chan<- []model.MetricInfoRaw
	metric MetricInfoRawer
}

// NewCollector - создание коллектора, сборщика метрик, где:
//   - in - при получении данных с данного канала запускается формирование метрик;
//   - m - функция формирование метрик;
//   - out - сформированные метрики отправляются в данный канал.
func NewCollector(in <-chan struct{}, m MetricInfoRawer) (*state, <-chan []model.MetricInfoRaw) {
	out := make(chan []model.MetricInfoRaw)
	return &state{
		in:     in,
		out:    out,
		metric: m,
	}, out
}

// Do - запуск сборщика метрик.
func (s *state) Do(ctx context.Context) {
	for {
		select {
		case <-s.in:
			s.out <- s.metric.MetricInfoRaw()
		case <-ctx.Done():
			log.Printf("Collecter closed with cause:%s\n", ctx.Err())
			return
		}
	}
}
