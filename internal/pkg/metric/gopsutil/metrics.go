// Package gopsutil for collect metrics from gopsutil package.
package gopsutil

import (
	"strconv"

	"github.com/k0st1a/metrics/internal/pkg/rawmetric"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type state struct {
}

// NewMetric - создание сущности по упаковки метрик из пакета gopsutil в формат rawmetric.MetricInfoRaw.
func NewMetric() *state {
	return &state{}
}

// MetricInfoRaw - упаковка метрик из пакета gopsutil в формат rawmetric.MetricInfoRaw.
func (s *state) MetricInfoRaw() []rawmetric.MetricInfoRaw {
	mi := s.mem2MetricInfoRaw()
	ci := s.cpu2MetricInfoRaw()
	return append(mi, ci...)
}

// mem2MetricInfoRaw - упаковка метрик `TotalMemory`, `FreeMemory` из пакета `github.com/shirou/gopsutil/v3/mem`
// в формат rawmetric.MetricInfoRaw.
func (s *state) mem2MetricInfoRaw() []rawmetric.MetricInfoRaw {
	mem, err := mem.VirtualMemory()
	if err != nil {
		log.Error().Err(err).Msg("get memory information error")
		return []rawmetric.MetricInfoRaw{}
	}

	return []rawmetric.MetricInfoRaw{
		rawmetric.MetricInfoRaw{
			Name:  "TotalMemory",
			Type:  rawmetric.Gauge,
			Value: float64(mem.Total),
		},
		rawmetric.MetricInfoRaw{
			Name:  "FreeMemory",
			Type:  rawmetric.Gauge,
			Value: float64(mem.Free),
		},
	}
}

// cpu2MetricInfoRaw - упаковка метрики `CPUutilization` из пакета `github.com/shirou/gopsutil/v3/cpu`
// в формат rawmetric.MetricInfoRaw.
func (s *state) cpu2MetricInfoRaw() []rawmetric.MetricInfoRaw {
	cpu, err := cpu.Percent(0, true)
	if err != nil {
		log.Error().Err(err).Msg("get cpu percent usage information error")
		return []rawmetric.MetricInfoRaw{}
	}

	mi := make([]rawmetric.MetricInfoRaw, len(cpu))

	for i, v := range cpu {
		mi[i] = rawmetric.MetricInfoRaw{
			Name:  "CPUutilization" + strconv.Itoa(i),
			Type:  rawmetric.Gauge,
			Value: float64(v),
		}
	}

	return mi
}
