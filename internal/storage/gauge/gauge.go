package gauge

import (
	"fmt"
	"strconv"
)

type Storage interface {
	GetGauge(string) (float64, bool)
	StoreGauge(string, float64)
}

func Store(name, value string, storage Storage) error {
	v, err := parser(value)
	if err != nil {
		return fmt.Errorf("gauge parse error:%w", err)
	}

	storage.StoreGauge(name, v)
	return nil
}

func Get(name string, storage Storage) (string, bool) {
	v, ok := storage.GetGauge(name)
	if !ok {
		return "", ok
	}

	return stringer(v), ok
}

func parser(s string) (float64, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return v, fmt.Errorf("parse float error:%w", err)
	}

	return v, nil
}

func stringer(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
