// Package ports contatins errors of metrics.
package ports

import "errors"

var (
	ErrMetricsNoCounter = errors.New("metrics: no counter")
	ErrMetricsNoGauge   = errors.New("metrics: no gauge")
)
