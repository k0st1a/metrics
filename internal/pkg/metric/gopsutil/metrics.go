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

// NewMetric - создание сущности по упаковки метрик из пакета gopsutil в формат rawmetric.Info.
func NewMetric() *state {
	return &state{}
}

// Info - упаковка метрик из пакета gopsutil в формат rawmetric.Info.
func (s *state) Info() []rawmetric.Info {
	mi := s.mem2Info()
	ci := s.cpu2Info()
	return append(mi, ci...)
}

// mem2Info - упаковка метрик `TotalMemory`, `FreeMemory` из пакета `github.com/shirou/gopsutil/v3/mem`
// в формат rawmetric.Info.
func (s *state) mem2Info() []rawmetric.Info {
	mem, err := mem.VirtualMemory()
	if err != nil {
		log.Error().Err(err).Msg("get memory information error")
		return []rawmetric.Info{}
	}

	return []rawmetric.Info{
		rawmetric.Info{
			Name:  "TotalMemory",
			Type:  rawmetric.Gauge,
			Value: float64(mem.Total),
		},
		rawmetric.Info{
			Name:  "FreeMemory",
			Type:  rawmetric.Gauge,
			Value: float64(mem.Free),
		},
	}
}

// cpu2Info - упаковка метрики `CPUutilization` из пакета `github.com/shirou/gopsutil/v3/cpu`
// в формат rawmetric.Info.
func (s *state) cpu2Info() []rawmetric.Info {
	cpu, err := cpu.Percent(0, true)
	if err != nil {
		log.Error().Err(err).Msg("get cpu percent usage information error")
		return []rawmetric.Info{}
	}

	mi := make([]rawmetric.Info, len(cpu))

	for i, v := range cpu {
		mi[i] = rawmetric.Info{
			Name:  "CPUutilization" + strconv.Itoa(i),
			Type:  rawmetric.Gauge,
			Value: float64(v),
		}
	}

	return mi
}
