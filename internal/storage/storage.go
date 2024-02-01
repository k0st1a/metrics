package storage

import (
	"github.com/k0st1a/metrics/internal/storage/file"
	"github.com/k0st1a/metrics/internal/storage/inmemory"
	"github.com/rs/zerolog/log"
)

type Storage interface {
	GetGauge(string) (float64, bool)
	StoreGauge(string, float64)

	GetCounter(string) (int64, bool)
	StoreCounter(string, int64)

	GetAll() (map[string]int64, map[string]float64)
}

func NewStorage(path string, interval int, restore bool) Storage {
	if path == "" {
		log.Debug().Msg("Using memory storage")
		return inmemory.NewStorage()
	}
	log.Debug().Msg("Using file storage")
	return file.NewStorage(path, interval, restore)
}
