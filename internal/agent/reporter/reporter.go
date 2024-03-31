package reporter

import (
	"time"

	"github.com/k0st1a/metrics/internal/metrics"
)

type Runner interface {
	Run()
}

type Metrics2MetricInfoer interface {
	Metrics2MetricInfo() []metrics.MetricInfo
}

type reporter struct {
	metrics        Metrics2MetricInfoer
	channel        chan<- []metrics.MetricInfo
	reportInterval int
}

func NewReporter(m Metrics2MetricInfoer, i int) (Runner, <-chan []metrics.MetricInfo) {
	mc := make(chan []metrics.MetricInfo)

	return &reporter{
		metrics:        m,
		channel:        mc,
		reportInterval: i,
	}, mc
}

func (r *reporter) Run() {
	ticker := time.NewTicker(time.Duration(r.reportInterval) * time.Second)

	for range ticker.C {
		r.channel <- r.metrics.Metrics2MetricInfo()
	}
}
