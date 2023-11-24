package counter

import (
	"strconv"
)

type Storage interface {
	GetCounter(string) (int64, bool)
	StoreCounter(string, int64)
}

func Store(name, value string, storage Storage) error {
	v, err := parser(value)
	if err != nil {
		return err
	}

	c, ok := storage.GetCounter(name)
	if ok {
		storage.StoreCounter(name, c+v)
		return nil
	}

	storage.StoreCounter(name, v)
	return nil
}

func Get(name string, storage Storage) (string, bool) {
	v, ok := storage.GetCounter(name)
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

	return v, err
}

func stringer(i int64) string {
	return strconv.FormatInt(i, 10)
}
