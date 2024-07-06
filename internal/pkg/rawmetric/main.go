// Package rawmetric for work with intermal representation of metrics.
package rawmetric

type Type int

const (
	Gauge = iota
	Counter
)

// Info - структура для хранения "сырых" метрик.
type Info struct {
	Value any
	Name  string
	Type  Type
}

// Append - добавление метрики в мапу.
func Append(acc map[string]Info, adding []Info) map[string]Info {
	for _, v := range adding {
		acc[v.Name] = Info{
			Name:  v.Name,
			Type:  v.Type,
			Value: v.Value,
		}
	}

	return acc
}

// Map2List  - преобразование из map в list.
func Map2List(m map[string]Info) []Info {
	l := make([]Info, len(m))
	i := 0
	for _, v := range m {
		l[i] = v
		i++
	}

	return l
}
