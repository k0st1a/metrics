package inmemorystorage

import "github.com/rs/zerolog/log"

type storage struct {
	gauge   map[string]float64
	counter map[string]int64
}

type Storage interface {
	StoreGauge(name string, value float64)
	GetGauge(name string) (value float64, ok bool)
	StoreCounter(name string, value int64)
	GetCounter(name string) (value int64, ok bool)
	GetAll() (counters map[string]int64, gauges map[string]float64)
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

func (s *storage) StoreGauge(name string, value float64) {
	log.Printf("StoreGauge, name(%v), value(%v)", name, value)
	s.gauge[name] = value
}

func (s *storage) GetGauge(name string) (float64, bool) {
	v, ok := s.gauge[name]
	log.Printf("GetGauge, name(%v), value(%v), ok(%v)", name, v, ok)
	return v, ok
}

func (s *storage) StoreCounter(name string, value int64) {
	log.Printf("StoreCounter, name(%v), value(%v)", name, value)
	s.counter[name] = value
}

func (s *storage) GetCounter(name string) (int64, bool) {
	v, ok := s.counter[name]
	log.Printf("GetCounter, name(%v), value(%v), ok(%v)", name, v, ok)
	return v, ok
}

func (s *storage) GetAll() (map[string]int64, map[string]float64) {
	return s.counter, s.gauge
}
