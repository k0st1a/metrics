package filestorage

import (
	"sync"

	"github.com/k0st1a/metrics/internal/storage/file/io"
	"github.com/k0st1a/metrics/internal/storage/inmemory"
	"github.com/rs/zerolog/log"
)

type Storage interface {
	GetGauge(name string) (bool float64, ok bool)
	StoreGauge(name string, value float64)

	GetCounter(name string) (value int64, ok bool)
	StoreCounter(name string, value int64)

	GetAll() (gauges map[string]int64, counters map[string]float64)
}

type Writer interface {
	Write(io.StorageGeter) error
}

type fileStorage struct {
	storage Storage
	writer  io.Writer
	mutex   sync.Mutex
}

func NewStorage(path string, interval int, restore bool) Storage {
	if path == "" {
		return inmemory.NewStorage()
	}

	var s Storage

	if restore {
		c, g, err := io.Read(path)
		if err != nil {
			log.Error().Err(err).Msg("io.Read Error")
		} else {
			s = inmemory.NewStorageWith(c, g)
		}
	}

	if s == nil {
		s = inmemory.NewStorage()
	}

	w := io.NewWriter(path)

	if interval != 0 {
		iw := io.NewIntervalWriter(w, s)
		go iw.Run(interval)
		w = nil
	}

	return &fileStorage{
		storage: s,
		writer:  w,
	}
}

func (s *fileStorage) StoreGauge(name string, value float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	log.Debug().
		Str("name:", name).
		Float64("value", value).
		Msg("StoreGauge")

	s.storage.StoreGauge(name, value)

	s.writeStorage()
}

func (s *fileStorage) GetGauge(name string) (float64, bool) {
	log.Debug().
		Str("name:", name).
		Msg("GetGauge")

	return s.storage.GetGauge(name)
}

func (s *fileStorage) StoreCounter(name string, value int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	log.Debug().
		Str("name", name).
		Int64("value", value).
		Msg("StoreCounter")

	s.storage.StoreCounter(name, value)

	s.writeStorage()
}

func (s *fileStorage) GetCounter(name string) (int64, bool) {
	log.Debug().
		Str("name", name).
		Msg("GetCounter")

	return s.storage.GetCounter(name)
}

func (s *fileStorage) GetAll() (map[string]int64, map[string]float64) {
	log.Debug().
		Msg("GetAll")

	return s.storage.GetAll()
}

func (s *fileStorage) writeStorage() {
	log.Debug().Msg("Write storage")
	if s.writer == nil {
		return
	}

	err := s.writer.Write(s.storage)
	if err != nil {
		log.Error().Err(err).Msg("write error storage to file")
	}
	log.Debug().Msg("Storage writed")
}
