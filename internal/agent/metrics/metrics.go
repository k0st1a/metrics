package metrics

import (
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

type MetricInfo struct {
	Type  string
	Value string
}

type MyStats struct {
	PollCount   uint64
	RandomValue float64
	MemStats    runtime.MemStats
}

func (metrics *MyStats) update() {
	runtime.ReadMemStats(&metrics.MemStats)
	metrics.RandomValue = rand.Float64()
}

func RunUpdateMetrics(metrics *MyStats, pollInternal int) {
	for {
		metrics.update()
		time.Sleep(time.Duration(pollInternal) * time.Second)
	}
}

func (metrics *MyStats) Prepare() map[string]MetricInfo {
	metrics.PollCount += 1

	return map[string]MetricInfo{
		"Alloc":         MetricInfo{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.Alloc, 10)},
		"BuckHashSys":   MetricInfo{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.BuckHashSys, 10)},
		"Frees":         MetricInfo{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.Frees, 10)},
		"GCSys":         MetricInfo{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.GCSys, 10)},
		"HeapAlloc":     MetricInfo{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.HeapAlloc, 10)},
		"HeapIdle":      MetricInfo{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.HeapIdle, 10)},
		"HeapInuse":     MetricInfo{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.HeapInuse, 10)},
		"HeapObjects":   MetricInfo{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.HeapObjects, 10)},
		"HeapReleased":  MetricInfo{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.HeapReleased, 10)},
		"HeapSys":       MetricInfo{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.HeapSys, 10)},
		"LastGC":        MetricInfo{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.LastGC, 10)},
		"Lookups":       MetricInfo{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.Lookups, 10)},
		"MCacheInuse":   MetricInfo{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.MCacheInuse, 10)},
		"MCacheSys":     MetricInfo{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.MCacheSys, 10)},
		"MSpanInuse":    MetricInfo{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.MSpanInuse, 10)},
		"MSpanSys":      MetricInfo{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.MSpanSys, 10)},
		"Mallocs":       MetricInfo{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.Mallocs, 10)},
		"NextGC":        MetricInfo{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.NextGC, 10)},
		"OtherSys":      MetricInfo{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.OtherSys, 10)},
		"PauseTotalNs":  MetricInfo{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.PauseTotalNs, 10)},
		"StackInuse":    MetricInfo{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.StackInuse, 10)},
		"StackSys":      MetricInfo{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.StackSys, 10)},
		"Sys":           MetricInfo{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.Sys, 10)},
		"TotalAlloc":    MetricInfo{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.TotalAlloc, 10)},
		"NumForcedGC":   MetricInfo{Type: "gauge", Value: strconv.FormatUint(uint64(metrics.MemStats.NumForcedGC), 10)},
		"NumGC":         MetricInfo{Type: "gauge", Value: strconv.FormatUint(uint64(metrics.MemStats.NumGC), 10)},
		"GCCPUFraction": MetricInfo{Type: "gauge", Value: strconv.FormatFloat(metrics.MemStats.GCCPUFraction, 'g', -1, 64)},

		"PollCount":   MetricInfo{Type: "counter", Value: strconv.FormatUint(metrics.PollCount, 10)},
		"RandomValue": MetricInfo{Type: "gauge", Value: strconv.FormatFloat(metrics.RandomValue, 'g', -1, 64)},
	}
}
