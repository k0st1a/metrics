package collector

import (
	"github.com/k0st1a/metrics/internal/agent/model"
)

type MetricInfoRawer interface {
	MetricInfoRaw() []model.MetricInfoRaw
}

type state struct {
	in     <-chan struct{}
	out    chan<- []model.MetricInfoRaw
	metric MetricInfoRawer
}

func NewCollector(in <-chan struct{}, m MetricInfoRawer) (*state, <-chan []model.MetricInfoRaw) {
	out := make(chan []model.MetricInfoRaw)
	return &state{
		in:     in,
		out:    out,
		metric: m,
	}, out
}

func (s *state) Do() {
	for range s.in {
		s.out <- s.metric.MetricInfoRaw()
	}
}
