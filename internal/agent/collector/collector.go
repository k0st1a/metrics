package collector

import (
	"github.com/k0st1a/metrics/internal/agent/model"
)

type MetricInfoer interface {
	MetricInfo() []model.MetricInfo
}

type Doer interface {
	Do()
}

type state struct {
	in     <-chan struct{}
	out    chan<- []model.MetricInfo
	metric MetricInfoer
}

func NewCollector(in <-chan struct{}, m MetricInfoer) (Doer, <-chan []model.MetricInfo) {
	out := make(chan []model.MetricInfo)
	return &state{
		in:     in,
		out:    out,
		metric: m,
	}, out
}

func (s *state) Do() {
	for range s.in {
		s.out <- s.metric.MetricInfo()
	}
}
