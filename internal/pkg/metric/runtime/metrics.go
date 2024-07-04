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
			Type:  "gauge",
			Value: s.memStats.Alloc,
		},
		model.MetricInfoRaw{
			Name:  "BuckHashSys",
			Type:  "gauge",
			Value: s.memStats.BuckHashSys,
		},
		model.MetricInfoRaw{
			Name:  "Frees",
			Type:  "gauge",
			Value: s.memStats.Frees,
		},
		model.MetricInfoRaw{
			Name:  "GCSys",
			Type:  "gauge",
			Value: s.memStats.GCSys,
		},
		model.MetricInfoRaw{
			Name:  "HeapAlloc",
			Type:  "gauge",
			Value: s.memStats.HeapAlloc,
		},
		model.MetricInfoRaw{
			Name:  "HeapIdle",
			Type:  "gauge",
			Value: s.memStats.HeapIdle,
		},
		model.MetricInfoRaw{
			Name:  "HeapInuse",
			Type:  "gauge",
			Value: s.memStats.HeapInuse,
		},
		model.MetricInfoRaw{
			Name:  "HeapObjects",
			Type:  "gauge",
			Value: s.memStats.HeapObjects,
		},
		model.MetricInfoRaw{
			Name:  "HeapReleased",
			Type:  "gauge",
			Value: s.memStats.HeapReleased,
		},
		model.MetricInfoRaw{
			Name:  "HeapSys",
			Type:  "gauge",
			Value: s.memStats.HeapSys,
		},
		model.MetricInfoRaw{
			Name:  "LastGC",
			Type:  "gauge",
			Value: s.memStats.LastGC,
		},
		model.MetricInfoRaw{
			Name:  "Lookups",
			Type:  "gauge",
			Value: s.memStats.Lookups,
		},
		model.MetricInfoRaw{
			Name:  "MCacheInuse",
			Type:  "gauge",
			Value: s.memStats.MCacheInuse,
		},
		model.MetricInfoRaw{
			Name:  "MCacheSys",
			Type:  "gauge",
			Value: s.memStats.MCacheSys,
		},
		model.MetricInfoRaw{
			Name:  "MSpanInuse",
			Type:  "gauge",
			Value: s.memStats.MSpanInuse,
		},
		model.MetricInfoRaw{
			Name:  "MSpanSys",
			Type:  "gauge",
			Value: s.memStats.MSpanSys,
		},
		model.MetricInfoRaw{
			Name:  "Mallocs",
			Type:  "gauge",
			Value: s.memStats.Mallocs,
		},
		model.MetricInfoRaw{
			Name:  "NextGC",
			Type:  "gauge",
			Value: s.memStats.NextGC,
		},
		model.MetricInfoRaw{
			Name:  "OtherSys",
			Type:  "gauge",
			Value: s.memStats.OtherSys,
		},
		model.MetricInfoRaw{
			Name:  "PauseTotalNs",
			Type:  "gauge",
			Value: s.memStats.PauseTotalNs,
		},
		model.MetricInfoRaw{
			Name:  "StackInuse",
			Type:  "gauge",
			Value: s.memStats.StackInuse,
		},
		model.MetricInfoRaw{
			Name:  "StackSys",
			Type:  "gauge",
			Value: s.memStats.StackSys,
		},
		model.MetricInfoRaw{
			Name:  "Sys",
			Type:  "gauge",
			Value: s.memStats.Sys,
		},
		model.MetricInfoRaw{
			Name:  "TotalAlloc",
			Type:  "gauge",
			Value: s.memStats.TotalAlloc,
		},
		model.MetricInfoRaw{
			Name:  "NumForcedGC",
			Type:  "gauge",
			Value: uint64(s.memStats.NumForcedGC),
		},
		model.MetricInfoRaw{
			Name:  "NumGC",
			Type:  "gauge",
			Value: uint64(s.memStats.NumGC),
		},
		model.MetricInfoRaw{
			Name:  "GCCPUFraction",
			Type:  "gauge",
			Value: s.memStats.GCCPUFraction,
		},
		model.MetricInfoRaw{
			Name:  "PollCount",
			Type:  "counter",
			Value: s.pollCount,
		},
		model.MetricInfoRaw{
			Name:  "RandomValue",
			Type:  "gauge",
			Value: s.randomValue,
		},
	}
}
