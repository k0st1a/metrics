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

func NewStorage(file string, interval int, restore bool) Storage {
	if file == "" {
		log.Debug().Msg("Using memory storage")
		return inmemorystorage.NewStorage()
	}
	log.Debug().Msg("Using file storage")
	return filestorage.NewStorage(file, interval, restore)
}
