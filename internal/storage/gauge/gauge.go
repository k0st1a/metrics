package gauge

import (
	"strconv"
	"strings"

	"github.com/k0st1a/metrics/internal/storage"
)

func Store(name, value string) error {
	v, err := parser(value)
	if err != nil {
		return err
	}

	storage.StoreGauge(name, v)
	return nil
}

func Get(name string) (string, bool) {
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
	return addDotIfNo(s)
}

func addDotIfNo(s string) string {
	if strings.ContainsRune(s, 46) { // 46 - ascii code of dot
		return s
	}
	return s + "."
}
