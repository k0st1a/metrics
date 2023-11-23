package gauge

import (
	"fmt"
	"strconv"

	"github.com/k0st1a/metrics/internal/utils"
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
	return strconv.ParseFloat(s, 64)
}

func stringer(f float64) string {
	s := strconv.FormatFloat(f, 'g', -1, 64)
	return utils.AddDotIfNo(s)
}
