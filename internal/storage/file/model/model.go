package model

import (
	"fmt"

	"github.com/mailru/easyjson"
	"github.com/rs/zerolog/log"
)

//go:generate easyjson -all model.go
type Metric struct {
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае gauge
	Name  string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // тип метрики, принимающий значение gauge или counter
}

type Metrics struct {
	List []Metric `json:"list"`
}

func Deserialize(b []byte) (map[string]int64, map[string]float64, error) {
	m := Metrics{}
	err := easyjson.Unmarshal(b, &m)
	if err != nil {
		return nil, nil, fmt.Errorf("easyjson.Unmarshal error:%w", err)
	}

	c := make(map[string]int64)
	g := make(map[string]float64)

	for _, v := range m.List {
		switch v.MType {
		case "counter":
			c[v.Name] = *v.Delta
		case "gauge":
			g[v.Name] = *v.Value
		default:
			log.Error().Msg("unknown MType")
		}
	}

	return c, g, nil
}

func Serialize(c map[string]int64, g map[string]float64) ([]byte, error) {
	m := []Metric{}

	for k, v := range c {
		v2 := v
		m = append(m, Metric{Name: k, MType: "counter", Delta: &v2})
	}

	for k, v2 := range g {
		v3 := v2
		m = append(m, Metric{Name: k, MType: "gauge", Value: &v3})
	}

	b, err := easyjson.Marshal(&Metrics{List: m})
	if err != nil {
		return nil, fmt.Errorf("easyjson.Marshal error:%w", err)
	}

	return b, nil
}
