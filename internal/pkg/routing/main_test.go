package routing

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseHost(t *testing.T) {
	tests := []struct {
		name        string
		dst         string
		err         string
		expectedSrc string
	}{
		{
			name: "check error of split host port",
			dst:  "localhost",
			err:  "split address error",
		},
		{
			name: "check error of parse ip",
			dst:  "localhost:9090",
			err:  "bad host address",
		},
		{
			name:        "check success parse host",
			dst:         "127.0.0.1:9090",
			expectedSrc: "127.0.0.1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			src, err := ParseHost(test.dst)
			if test.err != "" {
				assert.Nil(t, src)
				require.ErrorContains(t, err, test.err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.expectedSrc, src.String())
		})
	}
}

func TestRoute(t *testing.T) {
	tests := []struct {
		name        string
		dst         net.IP
		err         string
		expectedSrc string
	}{
		{
			name: "check route error",
			dst:  nil,
			err:  "route error",
		},
		{
			name:        "check route",
			dst:         net.ParseIP("127.0.0.1"),
			expectedSrc: "127.0.0.1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router, err := New()
			require.NoError(t, err)

			src, err := router.Route(test.dst)
			if test.err != "" {
				assert.Nil(t, src)
				require.ErrorContains(t, err, test.err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.expectedSrc, src.String())
		})
	}
}
