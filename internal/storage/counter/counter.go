package counter

import (
	"fmt"
	"strconv"
)

type Storage interface {
	GetCounter(string) (int64, bool)
	StoreCounter(string, int64)
}

type counterStorage struct {
	storage Storage
}

func NewCounterStorage(s Storage) counterStorage {
	return counterStorage{
		storage: s,
	}
}

func (s counterStorage) Store(name, value string) error {
	v, err := parser(value)
	if err != nil {
		return err
	}

	c, ok := s.storage.GetCounter(name)
	if ok {
		s.storage.StoreCounter(name, c+v)
		return nil
	}

	s.storage.StoreCounter(name, v)
	return nil
}

func (s counterStorage) Get(name string) (string, bool) {
	v, ok := s.storage.GetCounter(name)
	if !ok {
		return "", ok
	}

	return stringer(v), ok
}

func parser(s string) (int64, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return v, fmt.Errorf("parse int error:%w", err)
	}

	return v, nil
}

func stringer(i int64) string {
	return strconv.FormatInt(i, 10)
}
