package inmemory

import (
	"context"

	"github.com/rs/zerolog/log"
)

type storage struct {
	gauge   map[string]float64
	counter map[string]int64
}

type Storage interface {
	GetGauge(ctx context.Context, name string) (*float64, error)
	StoreGauge(ctx context.Context, name string, value float64) error

	GetCounter(ctx context.Context, name string) (*int64, error)
	StoreCounter(ctx context.Context, name string, value int64) error

	GetAll(ctx context.Context) (counter map[string]int64, gauge map[string]float64, err error)
}

func NewStorage() Storage {
	return &storage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func NewStorageWith(c map[string]int64, g map[string]float64) Storage {
	return &storage{
		counter: c,
		gauge:   g,
	}
}

func (s *storage) StoreGauge(ctx context.Context, name string, value float64) error {
	log.Printf("StoreGauge, name(%v), value(%v)", name, value)
	s.gauge[name] = value
	return nil
}

func (s *storage) GetGauge(ctx context.Context, name string) (*float64, error) {
	v, ok := s.gauge[name]
	log.Printf("GetGauge, name(%v), value(%v), ok(%v)", name, v, ok)
	if ok {
		return &v, nil
	}
	return nil, nil
}

func (s *storage) StoreCounter(ctx context.Context, name string, value int64) error {
	log.Printf("StoreCounter, name(%v), value(%v)", name, value)
	s.counter[name] = value
	return nil
}

func (s *storage) GetCounter(ctx context.Context, name string) (*int64, error) {
	v, ok := s.counter[name]
	log.Printf("GetCounter, name(%v), value(%v), ok(%v)", name, v, ok)
	if ok {
		return &v, nil
	}
	return nil, nil
}

func (s *storage) GetAll(ctx context.Context) (map[string]int64, map[string]float64, error) {
	return s.counter, s.gauge, nil
}
