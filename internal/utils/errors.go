package utils

import "errors"

var (
	ErrMetricsNoCounter = errors.New("metrics: no counter")
	ErrMetricsNoGauge   = errors.New("metrics: no gauge")
)
