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

// NewMetric - создание сущности по упаковки метрик из пакета runtime в формат rawmetric.MetricInfoRaw.
func NewMetric() *state {
	return &state{}
}

// MetricInfoRaw - упаковка метрик из пакета runtime в формат rawmetric.MetricInfoRaw.
func (s *state) MetricInfoRaw() []rawmetric.MetricInfoRaw {
	s.update()
	return s.mem2MetricInfoRaw()
}

// update - вычитывание метрик из пакета runtime + обновление поля randomValue и pollCount.
func (s *state) update() {
	runtime.ReadMemStats(&s.memStats)
	s.randomValue = rand.Float64()
	s.pollCount++
}

// mem2MetricInfoRaw - упаковка метрик из пакета runtime в формат rawmetric.MetricInfoRaw.
func (s *state) mem2MetricInfoRaw() []rawmetric.MetricInfoRaw {
	return []rawmetric.MetricInfoRaw{
		rawmetric.MetricInfoRaw{
			Name:  "Alloc",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.Alloc),
		},
		rawmetric.MetricInfoRaw{
			Name:  "BuckHashSys",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.BuckHashSys),
		},
		rawmetric.MetricInfoRaw{
			Name:  "Frees",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.Frees),
		},
		rawmetric.MetricInfoRaw{
			Name:  "GCSys",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.GCSys),
		},
		rawmetric.MetricInfoRaw{
			Name:  "HeapAlloc",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.HeapAlloc),
		},
		rawmetric.MetricInfoRaw{
			Name:  "HeapIdle",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.HeapIdle),
		},
		rawmetric.MetricInfoRaw{
			Name:  "HeapInuse",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.HeapInuse),
		},
		rawmetric.MetricInfoRaw{
			Name:  "HeapObjects",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.HeapObjects),
		},
		rawmetric.MetricInfoRaw{
			Name:  "HeapReleased",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.HeapReleased),
		},
		rawmetric.MetricInfoRaw{
			Name:  "HeapSys",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.HeapSys),
		},
		rawmetric.MetricInfoRaw{
			Name:  "LastGC",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.LastGC),
		},
		rawmetric.MetricInfoRaw{
			Name:  "Lookups",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.Lookups),
		},
		rawmetric.MetricInfoRaw{
			Name:  "MCacheInuse",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.MCacheInuse),
		},
		rawmetric.MetricInfoRaw{
			Name:  "MCacheSys",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.MCacheSys),
		},
		rawmetric.MetricInfoRaw{
			Name:  "MSpanInuse",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.MSpanInuse),
		},
		rawmetric.MetricInfoRaw{
			Name:  "MSpanSys",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.MSpanSys),
		},
		rawmetric.MetricInfoRaw{
			Name:  "Mallocs",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.Mallocs),
		},
		rawmetric.MetricInfoRaw{
			Name:  "NextGC",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.NextGC),
		},
		rawmetric.MetricInfoRaw{
			Name:  "OtherSys",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.OtherSys),
		},
		rawmetric.MetricInfoRaw{
			Name:  "PauseTotalNs",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.PauseTotalNs),
		},
		rawmetric.MetricInfoRaw{
			Name:  "StackInuse",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.StackInuse),
		},
		rawmetric.MetricInfoRaw{
			Name:  "StackSys",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.StackSys),
		},
		rawmetric.MetricInfoRaw{
			Name:  "Sys",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.Sys),
		},
		rawmetric.MetricInfoRaw{
			Name:  "TotalAlloc",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.TotalAlloc),
		},
		rawmetric.MetricInfoRaw{
			Name:  "NumForcedGC",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.NumForcedGC),
		},
		rawmetric.MetricInfoRaw{
			Name:  "NumGC",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.NumGC),
		},
		rawmetric.MetricInfoRaw{
			Name:  "GCCPUFraction",
			Type:  rawmetric.Gauge,
			Value: float64(s.memStats.GCCPUFraction),
		},
		rawmetric.MetricInfoRaw{
			Name:  "PollCount",
			Type:  rawmetric.Counter,
			Value: int64(s.pollCount),
		},
		rawmetric.MetricInfoRaw{
			Name:  "RandomValue",
			Type:  rawmetric.Gauge,
			Value: float64(s.randomValue),
		},
	}
}
