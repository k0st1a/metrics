package gauge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddDotIfNo(t *testing.T) {
	tests := []struct {
		name        string
		str         string
		expectedStr string
	}{
		{
			name:        "no dot in str",
			str:         "88888",
			expectedStr: "88888.",
		},
		{
			name:        "dot in str",
			str:         "88888.",
			expectedStr: "88888.",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedStr, addDotIfNo(test.str))
		})
	}
}
