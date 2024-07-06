// Package runtime for collect metrics from runtime package.
package runtime

import (
	"math/rand"
	"runtime"

	"github.com/k0st1a/metrics/internal/agent/model"
)

type state struct {
	pollCount   uint64
	randomValue float64
	memStats    runtime.MemStats
}

// NewMetric - создание сущности по упаковки метрик из пакета runtime в формат model.MetricInfoRaw.
func NewMetric() *state {
	return &state{}
}

// MetricInfoRaw - упаковка метрик из пакета runtime в формат model.MetricInfoRaw.
func (s *state) MetricInfoRaw() []model.MetricInfoRaw {
	s.update()
	return s.mem2MetricInfoRaw()
}

// update - вычитывание метрик из пакета runtime + обновление поля randomValue и pollCount.
func (s *state) update() {
	runtime.ReadMemStats(&s.memStats)
	s.randomValue = rand.Float64()
	s.pollCount++
}

// mem2MetricInfoRaw - упаковка метрик из пакета runtime в формат model.MetricInfoRaw.
func (s *state) mem2MetricInfoRaw() []model.MetricInfoRaw {
	return []model.MetricInfoRaw{
		model.MetricInfoRaw{
			Name:  "Alloc",
			Type:  model.Gauge,
			Value: float64(s.memStats.Alloc),
		},
		model.MetricInfoRaw{
			Name:  "BuckHashSys",
			Type:  model.Gauge,
			Value: float64(s.memStats.BuckHashSys),
		},
		model.MetricInfoRaw{
			Name:  "Frees",
			Type:  model.Gauge,
			Value: float64(s.memStats.Frees),
		},
		model.MetricInfoRaw{
			Name:  "GCSys",
			Type:  model.Gauge,
			Value: float64(s.memStats.GCSys),
		},
		model.MetricInfoRaw{
			Name:  "HeapAlloc",
			Type:  model.Gauge,
			Value: float64(s.memStats.HeapAlloc),
		},
		model.MetricInfoRaw{
			Name:  "HeapIdle",
			Type:  model.Gauge,
			Value: float64(s.memStats.HeapIdle),
		},
		model.MetricInfoRaw{
			Name:  "HeapInuse",
			Type:  model.Gauge,
			Value: float64(s.memStats.HeapInuse),
		},
		model.MetricInfoRaw{
			Name:  "HeapObjects",
			Type:  model.Gauge,
			Value: float64(s.memStats.HeapObjects),
		},
		model.MetricInfoRaw{
			Name:  "HeapReleased",
			Type:  model.Gauge,
			Value: float64(s.memStats.HeapReleased),
		},
		model.MetricInfoRaw{
			Name:  "HeapSys",
			Type:  model.Gauge,
			Value: float64(s.memStats.HeapSys),
		},
		model.MetricInfoRaw{
			Name:  "LastGC",
			Type:  model.Gauge,
			Value: float64(s.memStats.LastGC),
		},
		model.MetricInfoRaw{
			Name:  "Lookups",
			Type:  model.Gauge,
			Value: float64(s.memStats.Lookups),
		},
		model.MetricInfoRaw{
			Name:  "MCacheInuse",
			Type:  model.Gauge,
			Value: float64(s.memStats.MCacheInuse),
		},
		model.MetricInfoRaw{
			Name:  "MCacheSys",
			Type:  model.Gauge,
			Value: float64(s.memStats.MCacheSys),
		},
		model.MetricInfoRaw{
			Name:  "MSpanInuse",
			Type:  model.Gauge,
			Value: float64(s.memStats.MSpanInuse),
		},
		model.MetricInfoRaw{
			Name:  "MSpanSys",
			Type:  model.Gauge,
			Value: float64(s.memStats.MSpanSys),
		},
		model.MetricInfoRaw{
			Name:  "Mallocs",
			Type:  model.Gauge,
			Value: float64(s.memStats.Mallocs),
		},
		model.MetricInfoRaw{
			Name:  "NextGC",
			Type:  model.Gauge,
			Value: float64(s.memStats.NextGC),
		},
		model.MetricInfoRaw{
			Name:  "OtherSys",
			Type:  model.Gauge,
			Value: float64(s.memStats.OtherSys),
		},
		model.MetricInfoRaw{
			Name:  "PauseTotalNs",
			Type:  model.Gauge,
			Value: float64(s.memStats.PauseTotalNs),
		},
		model.MetricInfoRaw{
			Name:  "StackInuse",
			Type:  model.Gauge,
			Value: float64(s.memStats.StackInuse),
		},
		model.MetricInfoRaw{
			Name:  "StackSys",
			Type:  model.Gauge,
			Value: float64(s.memStats.StackSys),
		},
		model.MetricInfoRaw{
			Name:  "Sys",
			Type:  model.Gauge,
			Value: float64(s.memStats.Sys),
		},
		model.MetricInfoRaw{
			Name:  "TotalAlloc",
			Type:  model.Gauge,
			Value: float64(s.memStats.TotalAlloc),
		},
		model.MetricInfoRaw{
			Name:  "NumForcedGC",
			Type:  model.Gauge,
			Value: float64(s.memStats.NumForcedGC),
		},
		model.MetricInfoRaw{
			Name:  "NumGC",
			Type:  model.Gauge,
			Value: float64(s.memStats.NumGC),
		},
		model.MetricInfoRaw{
			Name:  "GCCPUFraction",
			Type:  model.Gauge,
			Value: float64(s.memStats.GCCPUFraction),
		},
		model.MetricInfoRaw{
			Name:  "PollCount",
			Type:  model.Counter,
			Value: int64(s.pollCount),
		},
		model.MetricInfoRaw{
			Name:  "RandomValue",
			Type:  model.Gauge,
			Value: float64(s.randomValue),
		},
	}
}
