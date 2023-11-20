package gauge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringer(t *testing.T) {
	tests := []struct {
		name        string
		value       float64
		expectedStr string
	}{
		{
			name:        "dot in value",
			value:       8888.1,
			expectedStr: "8888.1",
		},
		{
			name:        "no dot in value",
			value:       88888,
			expectedStr: "88888.",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedStr, stringer(test.value))
		})
	}
}
