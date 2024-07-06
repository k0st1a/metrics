package ports

import "github.com/k0st1a/metrics/internal/pkg/rawmetric"

type DoBatcher interface {
	DoBatch([]rawmetric.MetricInfoRaw) error
}

type Doer interface {
	Do(rawmetric.MetricInfoRaw) error
}
