// Package model for work with models of metrics.
package model

type Type int

const (
	Gauge = iota
	Counter
)

// MetricInfoRaw - структура для хранения "сырых" метрик.
type MetricInfoRaw struct {
	Value any
	Name  string
	Type  Type
}

// Append - добавление метрики в мапу.
func Append(acc map[string]MetricInfoRaw, adding []MetricInfoRaw) map[string]MetricInfoRaw {
	for _, v := range adding {
		acc[v.Name] = MetricInfoRaw{
			Name:  v.Name,
			Type:  v.Type,
			Value: v.Value,
		}
	}

	return acc
}

// Map2List  - преобразование из map в list.
func Map2List(m map[string]MetricInfoRaw) []MetricInfoRaw {
	l := make([]MetricInfoRaw, len(m))
	i := 0
	for _, v := range m {
		l[i] = v
		i++
	}

	return l
}
