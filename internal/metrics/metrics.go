package metrics

import (
	"math/rand"
	"runtime"
	"time"
)

type MyStats struct {
	PollCount   uint64
	RandomValue float64
	MemStats    runtime.MemStats
}

func RunUpdateMetrics(metrics *MyStats, pollInterval int) {
	ticker := time.NewTicker(time.Duration(pollInterval) * time.Second)

	for range ticker.C {
		metrics.update()
	}
}

func (metrics *MyStats) update() {
	runtime.ReadMemStats(&metrics.MemStats)
	metrics.RandomValue = rand.Float64()
}

func (metrics *MyStats) IncreasePollCount() {
	metrics.PollCount++
}
