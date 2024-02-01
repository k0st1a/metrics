package io

import (
	"fmt"
	"os"

	"github.com/k0st1a/metrics/internal/storage/file/model"
	"github.com/rs/zerolog/log"
)

type StorageGeter interface {
	GetAll() (counters map[string]int64, gauges map[string]float64)
}

type Writer interface {
	Write(StorageGeter) error
}

type file struct {
	path string
}

func NewWriter(p string) Writer {
	return &file{path: p}
}

func (f *file) Write(s StorageGeter) error {
	log.Printf("Write storage to file:%v", f.path)

	c, g := s.GetAll()

	p, err := model.Serialize(c, g)
	if err != nil {
		return fmt.Errorf("model.Serialize error:%w", err)
	}

	err = os.WriteFile(f.path, p, 0666)
	if err != nil {
		return fmt.Errorf("os.WriteFile error:%w", err)
	}

	log.Printf("Storage writed to file:%v", f.path)
	return nil
}

func Read(path string) (map[string]int64, map[string]float64, error) {
	log.Printf("Read storage from file:%v", path)

	p, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("os.ReadFile error:%w", err)
	}

	c, g, err := model.Deserialize(p)
	if err != nil {
		return nil, nil, fmt.Errorf("model.Deserialize error:%w", err)
	}

	log.Printf("Storage readed from file:%v", path)
	return c, g, nil
}
