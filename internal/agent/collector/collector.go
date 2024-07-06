// Package collector for collect metrics.
package collector

import (
	"context"

	"github.com/k0st1a/metrics/internal/pkg/rawmetric"
	"github.com/rs/zerolog/log"
)

// Infoer - интерфейс формирования метрик.
type Infoer interface {
	Info() []rawmetric.Info
}

type state struct {
	in     <-chan struct{}
	out    chan<- []rawmetric.Info
	metric Infoer
}

// NewCollector - создание коллектора, сборщика метрик, где:
//   - in - при получении данных с данного канала запускается формирование метрик;
//   - m - функция формирование метрик;
//   - out - сформированные метрики отправляются в данный канал.
func NewCollector(in <-chan struct{}, m Infoer) (*state, <-chan []rawmetric.Info) {
	out := make(chan []rawmetric.Info)
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
			s.out <- s.metric.Info()
		case <-ctx.Done():
			log.Printf("Collecter closed with cause:%s\n", ctx.Err())
			return
		}
	}
}
