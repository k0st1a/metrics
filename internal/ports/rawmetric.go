package ports

import "github.com/k0st1a/metrics/internal/pkg/rawmetric"

// RawMetricInfoer - интерфейс формирования "cырых" метрик.
type RawMetricInfoer interface {
	RawMetricInfo() []rawmetric.Info
}
