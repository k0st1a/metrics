// Package inmemory for save metrics to inmemory storage.
package inmemory

import (
	"context"
	"maps"

	"github.com/k0st1a/metrics/internal/utils"
	"github.com/rs/zerolog/log"
)

// Storage - внутреннее хранилище метрик.
type Storage struct {
	gauge   map[string]float64
	counter map[string]int64
}

// NewStorage - создать storage для хранения метрик в RAM.
func NewStorage() *Storage {
	return &Storage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

// NewStorageWith - создать storage для хранения метрик в RAM с заданными метриками типа counter и gauge, где:
//   - counter - метрики типа counter;
//   - gauge - метрики типа gauge.
func NewStorageWith(counter map[string]int64, gauge map[string]float64) *Storage {
	return &Storage{
		counter: counter,
		gauge:   gauge,
	}
}

// StoreGauge - сохраняет метрику типа gauge с именем name и значенем value.
func (s *Storage) StoreGauge(ctx context.Context, name string, value float64) error {
	log.Printf("StoreGauge, name(%v), value(%v)", name, value)
	s.gauge[name] = value
	return nil
}

// GetGauge - возвращает метрику типа gauge с именем name.
func (s *Storage) GetGauge(ctx context.Context, name string) (*float64, error) {
	v, ok := s.gauge[name]
	log.Printf("GetGauge, name(%v), value(%v), ok(%v)", name, v, ok)
	if ok {
		return &v, nil
	}

	return nil, utils.ErrMetricsNoGauge
}

// StoreCounter - сохраняет метрику типа counter с именем name и значенем value.
func (s *Storage) StoreCounter(ctx context.Context, name string, value int64) error {
	log.Printf("StoreCounter, name(%v), value(%v)", name, value)
	s.counter[name] += value
	return nil
}

// GetCounter - возвращает метрику типа gauge с именем name.
func (s *Storage) GetCounter(ctx context.Context, name string) (*int64, error) {
	v, ok := s.counter[name]
	log.Printf("GetCounter, name(%v), value(%v), ok(%v)", name, v, ok)
	if ok {
		return &v, nil
	}
	return nil, utils.ErrMetricsNoCounter
}

// StoreAll - сохраняет группу метрик типа counter и gauge.
func (s *Storage) StoreAll(ctx context.Context, counter map[string]int64, gauge map[string]float64) error {
	maps.Copy(s.counter, counter)
	maps.Copy(s.gauge, gauge)

	return nil
}

// GetAll - возвращает все метрики типа counter и gauge.
func (s *Storage) GetAll(ctx context.Context) (map[string]int64, map[string]float64, error) {
	return s.counter, s.gauge, nil
}
