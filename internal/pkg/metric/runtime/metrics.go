// Package runtime for collect metrics from runtime package.
package runtime

import (
	"math/rand"
	"runtime"

	"github.com/k0st1a/metrics/internal/pkg/rawmetric"
)

type state struct {
	pollCount   uint64
	randomValue float64
	memStats    runtime.MemStats
}

// NewMetric - создание сущности по упаковки метрик из пакета runtime в формат rawmetric.Info.
func NewMetric() *state {
	return &state{}
}

// RawMetricInfo - упаковка метрик из пакета runtime в формат rawmetric.Info.
func (s *state) RawMetricInfo() []rawmetric.Info {
	s.update()
	return s.mem2RawMetricInfo()
}

// update - вычитывание метрик из пакета runtime + обновление поля randomValue и pollCount.
func (s *state) update() {
	runtime.ReadMemStats(&s.memStats)
	s.randomValue = rand.Float64()
	s.pollCount++
}

// mem2RawMetricInfo - упаковка метрик из пакета runtime в формат rawmetric.Info.
func (s *state) mem2RawMetricInfo() []rawmetric.Info {
	return []rawmetric.Info{
		rawmetric.Info{
			Name:  "Alloc",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.Alloc),
		},
		rawmetric.Info{
			Name:  "BuckHashSys",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.BuckHashSys),
		},
		rawmetric.Info{
			Name:  "Frees",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.Frees),
		},
		rawmetric.Info{
			Name:  "GCSys",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.GCSys),
		},
		rawmetric.Info{
			Name:  "HeapAlloc",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.HeapAlloc),
		},
		rawmetric.Info{
			Name:  "HeapIdle",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.HeapIdle),
		},
		rawmetric.Info{
			Name:  "HeapInuse",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.HeapInuse),
		},
		rawmetric.Info{
			Name:  "HeapObjects",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.HeapObjects),
		},
		rawmetric.Info{
			Name:  "HeapReleased",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.HeapReleased),
		},
		rawmetric.Info{
			Name:  "HeapSys",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.HeapSys),
		},
		rawmetric.Info{
			Name:  "LastGC",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.LastGC),
		},
		rawmetric.Info{
			Name:  "Lookups",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.Lookups),
		},
		rawmetric.Info{
			Name:  "MCacheInuse",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.MCacheInuse),
		},
		rawmetric.Info{
			Name:  "MCacheSys",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.MCacheSys),
		},
		rawmetric.Info{
			Name:  "MSpanInuse",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.MSpanInuse),
		},
		rawmetric.Info{
			Name:  "MSpanSys",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.MSpanSys),
		},
		rawmetric.Info{
			Name:  "Mallocs",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.Mallocs),
		},
		rawmetric.Info{
			Name:  "NextGC",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.NextGC),
		},
		rawmetric.Info{
			Name:  "OtherSys",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.OtherSys),
		},
		rawmetric.Info{
			Name:  "PauseTotalNs",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.PauseTotalNs),
		},
		rawmetric.Info{
			Name:  "StackInuse",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.StackInuse),
		},
		rawmetric.Info{
			Name:  "StackSys",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.StackSys),
		},
		rawmetric.Info{
			Name:  "Sys",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.Sys),
		},
		rawmetric.Info{
			Name:  "TotalAlloc",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.TotalAlloc),
		},
		rawmetric.Info{
			Name:  "NumForcedGC",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.NumForcedGC),
		},
		rawmetric.Info{
			Name:  "NumGC",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.NumGC),
		},
		rawmetric.Info{
			Name:  "GCCPUFraction",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.GCCPUFraction),
		},
		rawmetric.Info{
			Name:  "PollCount",
			Type:  rawmetric.Counter,
			Value: int64(s.pollCount),
		},
		rawmetric.Info{
			Name:  "RandomValue",
			Type:  rawmetric.Gauge,
			Value: float64(s.randomValue),
		},
	}
}
