package metrics

import (
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

type MyStats struct {
	PollCount   uint64
	RandomValue float64
	MemStats    runtime.MemStats
}

type MetricInfo struct {
	Name  string
	MType string
	Value string
}

func RunUpdateMetrics(m *MyStats, pollInterval int) {
	ticker := time.NewTicker(time.Duration(pollInterval) * time.Second)

	for range ticker.C {
		m.update()
	}
}

func (m *MyStats) update() {
	runtime.ReadMemStats(&m.MemStats)
	m.RandomValue = rand.Float64()
}

func (m *MyStats) IncreasePollCount() {
	m.PollCount++
}

func (m *MyStats) Metrics2MetricInfo() []MetricInfo {
	return []MetricInfo{
		MetricInfo{
			Name:  "Alloc",
			MType: "gauge",
			Value: strconv.FormatUint(m.MemStats.Alloc, 10),
		},
		MetricInfo{
			Name:  "BuckHashSys",
			MType: "gauge",
			Value: strconv.FormatUint(m.MemStats.BuckHashSys, 10),
		},
		MetricInfo{
			Name:  "Frees",
			MType: "gauge",
			Value: strconv.FormatUint(m.MemStats.Frees, 10),
		},
		MetricInfo{
			Name:  "GCSys",
			MType: "gauge",
			Value: strconv.FormatUint(m.MemStats.GCSys, 10),
		},
		MetricInfo{
			Name:  "HeapAlloc",
			MType: "gauge",
			Value: strconv.FormatUint(m.MemStats.HeapAlloc, 10),
		},
		MetricInfo{
			Name:  "HeapIdle",
			MType: "gauge",
			Value: strconv.FormatUint(m.MemStats.HeapIdle, 10),
		},
		MetricInfo{
			Name:  "HeapInuse",
			MType: "gauge",
			Value: strconv.FormatUint(m.MemStats.HeapInuse, 10),
		},
		MetricInfo{
			Name:  "HeapObjects",
			MType: "gauge",
			Value: strconv.FormatUint(m.MemStats.HeapObjects, 10),
		},
		MetricInfo{
			Name:  "HeapReleased",
			MType: "gauge",
			Value: strconv.FormatUint(m.MemStats.HeapReleased, 10),
		},
		MetricInfo{
			Name:  "HeapSys",
			MType: "gauge",
			Value: strconv.FormatUint(m.MemStats.HeapSys, 10),
		},
		MetricInfo{
			Name:  "LastGC",
			MType: "gauge",
			Value: strconv.FormatUint(m.MemStats.LastGC, 10),
		},
		MetricInfo{
			Name:  "Lookups",
			MType: "gauge",
			Value: strconv.FormatUint(m.MemStats.Lookups, 10),
		},
		MetricInfo{
			Name:  "MCacheInuse",
			MType: "gauge",
			Value: strconv.FormatUint(m.MemStats.MCacheInuse, 10),
		},
		MetricInfo{
			Name:  "MCacheSys",
			MType: "gauge",
			Value: strconv.FormatUint(m.MemStats.MCacheSys, 10),
		},
		MetricInfo{
			Name:  "MSpanInuse",
			MType: "gauge",
			Value: strconv.FormatUint(m.MemStats.MSpanInuse, 10),
		},
		MetricInfo{
			Name:  "MSpanSys",
			MType: "gauge",
			Value: strconv.FormatUint(m.MemStats.MSpanSys, 10),
		},
		MetricInfo{
			Name:  "Mallocs",
			MType: "gauge",
			Value: strconv.FormatUint(m.MemStats.Mallocs, 10),
		},
		MetricInfo{
			Name:  "NextGC",
			MType: "gauge",
			Value: strconv.FormatUint(m.MemStats.NextGC, 10),
		},
		MetricInfo{
			Name:  "OtherSys",
			MType: "gauge",
			Value: strconv.FormatUint(m.MemStats.OtherSys, 10),
		},
		MetricInfo{
			Name:  "PauseTotalNs",
			MType: "gauge",
			Value: strconv.FormatUint(m.MemStats.PauseTotalNs, 10),
		},
		MetricInfo{
			Name:  "StackInuse",
			MType: "gauge",
			Value: strconv.FormatUint(m.MemStats.StackInuse, 10),
		},
		MetricInfo{
			Name:  "StackSys",
			MType: "gauge",
			Value: strconv.FormatUint(m.MemStats.StackSys, 10),
		},
		MetricInfo{
			Name:  "Sys",
			MType: "gauge",
			Value: strconv.FormatUint(m.MemStats.Sys, 10),
		},
		MetricInfo{
			Name:  "TotalAlloc",
			MType: "gauge",
			Value: strconv.FormatUint(m.MemStats.TotalAlloc, 10),
		},
		MetricInfo{
			Name:  "NumForcedGC",
			MType: "gauge",
			Value: strconv.FormatUint(uint64(m.MemStats.NumForcedGC), 10),
		},
		MetricInfo{
			Name:  "NumGC",
			MType: "gauge",
			Value: strconv.FormatUint(uint64(m.MemStats.NumGC), 10),
		},
		MetricInfo{
			Name:  "GCCPUFraction",
			MType: "gauge",
			Value: strconv.FormatFloat(m.MemStats.GCCPUFraction, 'g', -1, 64),
		},
		MetricInfo{
			Name:  "PollCount",
			MType: "counter",
			Value: strconv.FormatUint(m.PollCount, 10),
		},
		MetricInfo{
			Name:  "RandomValue",
			MType: "gauge",
			Value: strconv.FormatFloat(m.RandomValue, 'g', -1, 64),
		},
	}
}
