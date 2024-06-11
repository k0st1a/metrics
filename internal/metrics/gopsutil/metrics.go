package gopsutil

import (
	"strconv"

	"github.com/k0st1a/metrics/internal/agent/model"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type state struct {
}

// NewMetric - создание сущности по упаковки метрик из пакета gopsutil в формат model.MetricInfoRaw.
func NewMetric() *state {
	return &state{}
}

// MetricInfoRaw - упаковка метрик из пакета gopsutil в формат model.MetricInfoRaw.
func (s *state) MetricInfoRaw() []model.MetricInfoRaw {
	mi := s.mem2MetricInfoRaw()
	ci := s.cpu2MetricInfoRaw()
	return append(mi, ci...)
}

// mem2MetricInfoRaw - упаковка метрик `TotalMemory`, `FreeMemory` из пакета `github.com/shirou/gopsutil/v3/mem`
// в формат model.MetricInfoRaw.
func (s *state) mem2MetricInfoRaw() []model.MetricInfoRaw {
	mem, err := mem.VirtualMemory()
	if err != nil {
		log.Error().Err(err).Msg("get memory information error")
		return []model.MetricInfoRaw{}
	}

	return []model.MetricInfoRaw{
		model.MetricInfoRaw{
			Name:  "TotalMemory",
			Type:  "gauge",
			Value: mem.Total,
		},
		model.MetricInfoRaw{
			Name:  "FreeMemory",
			Type:  "gauge",
			Value: mem.Free,
		},
	}
}

// cpu2MetricInfoRaw - упаковка метрики `CPUutilization` из пакета `github.com/shirou/gopsutil/v3/cpu`
// в формат model.MetricInfoRaw.
func (s *state) cpu2MetricInfoRaw() []model.MetricInfoRaw {
	cpu, err := cpu.Percent(0, true)
	if err != nil {
		log.Error().Err(err).Msg("get cpu percent usage information error")
		return []model.MetricInfoRaw{}
	}

	mi := make([]model.MetricInfoRaw, len(cpu))

	for i, v := range cpu {
		mi[i] = model.MetricInfoRaw{
			Name:  "CPUutilization" + strconv.Itoa(i),
			Type:  "gauge",
			Value: v,
		}
	}

	return mi
}
