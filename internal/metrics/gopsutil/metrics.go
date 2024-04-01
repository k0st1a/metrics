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

func NewMetric() *state {
	return &state{}
}

func (s *state) MetricInfo() []model.MetricInfo {
	mi := s.mem2MetricInfo()
	ci := s.cpu2MetricInfo()
	return append(mi, ci...)
}

func (s *state) mem2MetricInfo() []model.MetricInfo {
	mem, err := mem.VirtualMemory()
	if err != nil {
		log.Error().Err(err).Msg("get memory information error")
		return []model.MetricInfo{}
	}

	return []model.MetricInfo{
		model.MetricInfo{
			Name:  "TotalMemory",
			MType: "gauge",
			Value: strconv.FormatUint(mem.Total, 10),
		},
		model.MetricInfo{
			Name:  "FreeMemory",
			MType: "gauge",
			Value: strconv.FormatUint(mem.Free, 10),
		},
	}
}

func (s *state) cpu2MetricInfo() []model.MetricInfo {
	cpu, err := cpu.Percent(0, true)
	if err != nil {
		log.Error().Err(err).Msg("get cpu percent usage information error")
		return []model.MetricInfo{}
	}

	mi := make([]model.MetricInfo, len(cpu))

	for i, v := range cpu {
		mi[i] = model.MetricInfo{
			Name:  "CPUutilization" + strconv.Itoa(i),
			MType: "gauge",
			Value: strconv.FormatFloat(v, 'g', -1, 64),
		}
	}

	return mi
}
