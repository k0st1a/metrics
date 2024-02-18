package io

import (
	"context"
	"fmt"
	"os"

	"github.com/k0st1a/metrics/internal/storage/file/model"
	"github.com/rs/zerolog/log"
)

type StorageGeter interface {
	GetAll(ctx context.Context) (counter map[string]int64, gauge map[string]float64, err error)
}

type Writer interface {
	Write(context.Context, StorageGeter) error
}

type file struct {
	path string
}

func NewWriter(p string) Writer {
	return &file{path: p}
}

const FileMode = 0600

func (f *file) Write(ctx context.Context, s StorageGeter) error {
	log.Printf("Write storage to file:%v", f.path)

	c, g, err := s.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("get all error:%w", err)
	}

	p, err := model.Serialize(c, g)
	if err != nil {
		return fmt.Errorf("model.Serialize error:%w", err)
	}

	err = os.WriteFile(f.path, p, FileMode)
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
