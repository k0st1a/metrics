package storage

import "github.com/rs/zerolog/log"

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (s *MemStorage) StoreGauge(name string, value float64) {
	log.Debug().
		Str("name:", name).
		Float64("value", value).
		Msg("StoreGauge")

	s.gauge[name] = value
}

func (s *MemStorage) GetGauge(name string) (float64, bool) {
	v, ok := s.gauge[name]
	log.Debug().
		Str("name:", name).
		Float64("value:", v).
		Msg("GetGauge")

	return v, ok
}

func (s *MemStorage) StoreCounter(name string, value int64) {
	log.Debug().
		Str("name", name).
		Int64("value", value).
		Msg("StoreCounter")

	s.counter[name] = value
}

func (s *MemStorage) GetCounter(name string) (int64, bool) {
	v, ok := s.counter[name]
	log.Debug().
		Str("name", name).
		Int64("value", v).
		Msg("GetCounter")

	return v, ok
}
