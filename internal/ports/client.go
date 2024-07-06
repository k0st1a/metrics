package ports

import "github.com/k0st1a/metrics/internal/agent/model"

type DoBatcher interface {
	DoBatch([]model.MetricInfoRaw) error
}

type Doer interface {
	Do(model.MetricInfoRaw) error
}
