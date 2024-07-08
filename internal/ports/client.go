package ports

import "github.com/k0st1a/metrics/internal/pkg/rawmetric"

type DoBatcher interface {
	DoBatch([]rawmetric.Info) error
}

type Doer interface {
	Do(rawmetric.Info) error
}
