package runtime

import (
	"math/rand"
	"runtime"
	"strconv"

	"github.com/k0st1a/metrics/internal/agent/model"
)

type state struct {
	pollCount   uint64
	randomValue float64
	memStats    runtime.MemStats
}

func NewMetric() *state {
	return &state{}
}

func (s *state) MetricInfo() []model.MetricInfo {
	s.update()
	return s.mem2MetricInfo()
}

func (s *state) update() {
	runtime.ReadMemStats(&s.memStats)
	s.randomValue = rand.Float64()
	s.pollCount++
}

func (s *state) mem2MetricInfo() []model.MetricInfo {
	return []model.MetricInfo{
		model.MetricInfo{
			Name:  "Alloc",
			MType: "gauge",
			Value: strconv.FormatUint(s.memStats.Alloc, 10),
		},
		model.MetricInfo{
			Name:  "BuckHashSys",
			MType: "gauge",
			Value: strconv.FormatUint(s.memStats.BuckHashSys, 10),
		},
		model.MetricInfo{
			Name:  "Frees",
			MType: "gauge",
			Value: strconv.FormatUint(s.memStats.Frees, 10),
		},
		model.MetricInfo{
			Name:  "GCSys",
			MType: "gauge",
			Value: strconv.FormatUint(s.memStats.GCSys, 10),
		},
		model.MetricInfo{
			Name:  "HeapAlloc",
			MType: "gauge",
			Value: strconv.FormatUint(s.memStats.HeapAlloc, 10),
		},
		model.MetricInfo{
			Name:  "HeapIdle",
			MType: "gauge",
			Value: strconv.FormatUint(s.memStats.HeapIdle, 10),
		},
		model.MetricInfo{
			Name:  "HeapInuse",
			MType: "gauge",
			Value: strconv.FormatUint(s.memStats.HeapInuse, 10),
		},
		model.MetricInfo{
			Name:  "HeapObjects",
			MType: "gauge",
			Value: strconv.FormatUint(s.memStats.HeapObjects, 10),
		},
		model.MetricInfo{
			Name:  "HeapReleased",
			MType: "gauge",
			Value: strconv.FormatUint(s.memStats.HeapReleased, 10),
		},
		model.MetricInfo{
			Name:  "HeapSys",
			MType: "gauge",
			Value: strconv.FormatUint(s.memStats.HeapSys, 10),
		},
		model.MetricInfo{
			Name:  "LastGC",
			MType: "gauge",
			Value: strconv.FormatUint(s.memStats.LastGC, 10),
		},
		model.MetricInfo{
			Name:  "Lookups",
			MType: "gauge",
			Value: strconv.FormatUint(s.memStats.Lookups, 10),
		},
		model.MetricInfo{
			Name:  "MCacheInuse",
			MType: "gauge",
			Value: strconv.FormatUint(s.memStats.MCacheInuse, 10),
		},
		model.MetricInfo{
			Name:  "MCacheSys",
			MType: "gauge",
			Value: strconv.FormatUint(s.memStats.MCacheSys, 10),
		},
		model.MetricInfo{
			Name:  "MSpanInuse",
			MType: "gauge",
			Value: strconv.FormatUint(s.memStats.MSpanInuse, 10),
		},
		model.MetricInfo{
			Name:  "MSpanSys",
			MType: "gauge",
			Value: strconv.FormatUint(s.memStats.MSpanSys, 10),
		},
		model.MetricInfo{
			Name:  "Mallocs",
			MType: "gauge",
			Value: strconv.FormatUint(s.memStats.Mallocs, 10),
		},
		model.MetricInfo{
			Name:  "NextGC",
			MType: "gauge",
			Value: strconv.FormatUint(s.memStats.NextGC, 10),
		},
		model.MetricInfo{
			Name:  "OtherSys",
			MType: "gauge",
			Value: strconv.FormatUint(s.memStats.OtherSys, 10),
		},
		model.MetricInfo{
			Name:  "PauseTotalNs",
			MType: "gauge",
			Value: strconv.FormatUint(s.memStats.PauseTotalNs, 10),
		},
		model.MetricInfo{
			Name:  "StackInuse",
			MType: "gauge",
			Value: strconv.FormatUint(s.memStats.StackInuse, 10),
		},
		model.MetricInfo{
			Name:  "StackSys",
			MType: "gauge",
			Value: strconv.FormatUint(s.memStats.StackSys, 10),
		},
		model.MetricInfo{
			Name:  "Sys",
			MType: "gauge",
			Value: strconv.FormatUint(s.memStats.Sys, 10),
		},
		model.MetricInfo{
			Name:  "TotalAlloc",
			MType: "gauge",
			Value: strconv.FormatUint(s.memStats.TotalAlloc, 10),
		},
		model.MetricInfo{
			Name:  "NumForcedGC",
			MType: "gauge",
			Value: strconv.FormatUint(uint64(s.memStats.NumForcedGC), 10),
		},
		model.MetricInfo{
			Name:  "NumGC",
			MType: "gauge",
			Value: strconv.FormatUint(uint64(s.memStats.NumGC), 10),
		},
		model.MetricInfo{
			Name:  "GCCPUFraction",
			MType: "gauge",
			Value: strconv.FormatFloat(s.memStats.GCCPUFraction, 'g', -1, 64),
		},
		model.MetricInfo{
			Name:  "PollCount",
			MType: "counter",
			Value: strconv.FormatUint(s.pollCount, 10),
		},
		model.MetricInfo{
			Name:  "RandomValue",
			MType: "gauge",
			Value: strconv.FormatFloat(s.randomValue, 'g', -1, 64),
		},
	}
}
