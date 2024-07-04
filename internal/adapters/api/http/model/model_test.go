package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		metrics []Metrics
	}{
		{
			name: "Serialize/Deserialize JSON to MetricsList",
			json: `[
					{
						"id": "CounterBatchZip215",
						"type": "counter",
						"delta": 2039207723
					},
					{
						"id": "GaugeBatchZip99",
						"type": "gauge",
						"value": 811072.5191439573
					},
					{
						"id": "CounterBatchZip215",
						"type": "counter",
						"delta": 329725296
					},
					{
						"id": "GaugeBatchZip99",
						"type": "gauge",
						"value": 723788.3421958535
					}
				   ]`,
			metrics: []Metrics{
				Metrics{
					ID:    "CounterBatchZip215",
					MType: "counter",
					Delta: adrInt64(2039207723),
				},
				Metrics{
					ID:    "GaugeBatchZip99",
					MType: "gauge",
					Value: adrFloat64(811072.5191439573),
				},
				Metrics{
					ID:    "CounterBatchZip215",
					MType: "counter",
					Delta: adrInt64(329725296),
				},
				Metrics{
					ID:    "GaugeBatchZip99",
					MType: "gauge",
					Value: adrFloat64(723788.3421958535),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m, err := DeserializeList([]byte(test.json))
			assert.NoError(t, err)
			assert.Equal(t, test.metrics, m)

			b, err := SerializeList(test.metrics)
			assert.NoError(t, err)
			assert.JSONEq(t, test.json, string(b))
		})
	}
}

func adrInt64(v int64) *int64 {
	return &v
}

func adrFloat64(v float64) *float64 {
	return &v
}
