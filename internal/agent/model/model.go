package model

import "strconv"

// MetricInfo - структура для "промежуточного" хранения метрик.
type MetricInfo struct {
	Name  string
	MType string
	Value string
}

// MetricInfoRaw - структура для хранения "сырых" метрик.
type MetricInfoRaw struct {
	Value any
	Name  string
	Type  string
}

// Append - добавление метрики в список.
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

// RawMap2InfoList  - преобразование метрик формата map в метрики формата list.
func RawMap2InfoList(r map[string]MetricInfoRaw) []MetricInfo {
	mi := make([]MetricInfo, len(r))
	i := 0
	for _, v := range r {
		mi[i] = raw2Info(v)
		i++
	}

	return mi
}

// raw2Info - преобразование "сырой" метрики в "промежуточный".
func raw2Info(m MetricInfoRaw) MetricInfo {
	var value string

	switch v := m.Value.(type) {
	case uint64:
		value = strconv.FormatUint(v, 10)
	case float64:
		value = strconv.FormatFloat(v, 'g', -1, 64)
	}

	return MetricInfo{
		Name:  m.Name,
		MType: m.Type,
		Value: value,
	}
}
