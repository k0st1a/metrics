package inmemory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemory(t *testing.T) {
	tests := []struct {
		name    string
		counter map[string]int64
		gauge   map[string]float64
	}{
		{
			name: "check NewStorageWith with counter and gauge",
			counter: map[string]int64{
				"counter1": 123,
				"counter2": 456,
			},
			gauge: map[string]float64{
				"gauge1": 123.1,
				"gauge2": 456.2,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := NewStorageWith(test.counter, test.gauge)

			c, g, err := s.GetAll(context.Background())
			assert.NoError(t, err)

			assert.Equal(t, test.counter, c)
			assert.Equal(t, test.gauge, g)
		})
	}
}
