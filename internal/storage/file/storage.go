package file

import (
	"context"
	"fmt"
	"sync"

	"github.com/k0st1a/metrics/internal/storage/file/io"
	"github.com/k0st1a/metrics/internal/storage/inmemory"
	"github.com/rs/zerolog/log"
)

type Storage interface {
	GetGauge(ctx context.Context, name string) (*float64, error)
	StoreGauge(ctx context.Context, name string, value float64) error

	GetCounter(ctx context.Context, name string) (*int64, error)
	StoreCounter(ctx context.Context, name string, value int64) error

	StoreAll(ctx context.Context, counter map[string]int64, gauge map[string]float64) error
	GetAll(ctx context.Context) (gauge map[string]int64, counter map[string]float64, err error)
}

type Writer interface {
	Write(io.StorageGeter) error
}

type fileStorage struct {
	storage Storage
	writer  io.Writer
	mutex   sync.Mutex
}

func NewStorage(ctx context.Context, path string, interval int, restore bool) Storage {
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
		go iw.Run(ctx, interval)
		w = nil
	}

	return &fileStorage{
		storage: s,
		writer:  w,
	}
}

func (s *fileStorage) StoreGauge(ctx context.Context, name string, value float64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	log.Debug().
		Str("name:", name).
		Float64("value", value).
		Msg("StoreGauge")

	err := s.storage.StoreGauge(ctx, name, value)
	if err != nil {
		return fmt.Errorf("store gauge error:%w", err)
	}

	s.writeStorage(ctx)

	return nil
}

func (s *fileStorage) GetGauge(ctx context.Context, name string) (*float64, error) {
	log.Debug().
		Str("name:", name).
		Msg("GetGauge")

	return s.storage.GetGauge(ctx, name)
}

func (s *fileStorage) StoreCounter(ctx context.Context, name string, value int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	log.Debug().
		Str("name", name).
		Int64("value", value).
		Msg("StoreCounter")

	err := s.storage.StoreCounter(ctx, name, value)
	if err != nil {
		return fmt.Errorf("store counter error:%w", err)
	}

	s.writeStorage(ctx)

	return nil
}

func (s *fileStorage) GetCounter(ctx context.Context, name string) (*int64, error) {
	log.Debug().
		Str("name", name).
		Msg("GetCounter")

	return s.storage.GetCounter(ctx, name)
}

func (s *fileStorage) StoreAll(ctx context.Context, counter map[string]int64, gauge map[string]float64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	log.Debug().
		Msg("StoreAll")

	err := s.storage.StoreAll(ctx, counter, gauge)
	if err != nil {
		return fmt.Errorf("store all error:%w", err)
	}

	s.writeStorage(ctx)

	return nil
}

func (s *fileStorage) GetAll(ctx context.Context) (map[string]int64, map[string]float64, error) {
	log.Debug().
		Msg("GetAll")

	return s.storage.GetAll(ctx)
}

func (s *fileStorage) writeStorage(ctx context.Context) {
	log.Debug().Msg("Write storage")
	if s.writer == nil {
		return
	}

	err := s.writer.Write(ctx, s.storage)
	if err != nil {
		log.Error().Err(err).Msg("write error storage to file")
	}
	log.Debug().Msg("Storage writed")
}
