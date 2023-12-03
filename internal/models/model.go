package models

import (
	"fmt"

	"github.com/mailru/easyjson"
)

//go:generate easyjson -all model.go
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func Deserialize(b []byte) (*Metrics, error) {
	m := &Metrics{}
	err := easyjson.Unmarshal(b, m)
	if err != nil {
		return nil, fmt.Errorf("easyjson.Unmarshal error:%w", err)
	}

	return m, nil
}

func Serialize(m *Metrics) ([]byte, error) {
	b, err := easyjson.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("easyjson.Marshal error:%w", err)
	}

	return b, nil
}
