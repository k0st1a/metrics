package counter

import (
	"strconv"

	"github.com/k0st1a/metrics/internal/storage"
)

func Store(name, value string) error {
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

func Get(name string) (string, bool) {
	v, ok := storage.GetCounter(name)
	if !ok {
		return "", ok
	}

	return stringer(v), ok
}

func parser(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func stringer(i int64) string {
	return strconv.FormatInt(i, 10)
}
