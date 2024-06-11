package models

import (
	"fmt"

	"github.com/mailru/easyjson"
)

//go:generate easyjson -all model.go

// Metrics - описание метрики для взаимодействия агента и сервера
//
//easyjson:json
type Metrics struct {
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
}

//easyjson:json
type MetricsList []Metrics

// Deserialize - распаковка байт в формат Metrics.
func Deserialize(b []byte) (*Metrics, error) {
	m := &Metrics{}
	err := easyjson.Unmarshal(b, m)
	if err != nil {
		return nil, fmt.Errorf("easyjson.Unmarshal error:%w", err)
	}

	return m, nil
}

// Serialize - упаковка Metrics в байты.
func Serialize(m *Metrics) ([]byte, error) {
	b, err := easyjson.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("easyjson.Marshal error:%w", err)
	}

	return b, nil
}

// DeserializeList - распаковка байт в формат []Metrics.
func DeserializeList(b []byte) ([]Metrics, error) {
	v := MetricsList{}
	err := easyjson.Unmarshal(b, &v)
	if err != nil {
		return nil, fmt.Errorf("easyjson.Unmarshal error:%w", err)
	}

	return []Metrics(v), nil
}

// SerializeList - упаковка []Metrics в байты.
func SerializeList(ml []Metrics) ([]byte, error) {
	v := MetricsList(ml)
	b, err := easyjson.Marshal(&v)
	if err != nil {
		return nil, fmt.Errorf("easyjson.Marshal error:%w", err)
	}

	return b, nil
}
