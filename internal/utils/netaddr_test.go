package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetAddressString(t *testing.T) {
	tests := []struct {
		name        string
		addr        NetAddress
		expectedStr string
	}{
		{
			name: "{host:\"localhost\", port:8080,} => localhost:8080",
			addr: NetAddress{
				host: "localhost",
				port: 8080,
			},
			expectedStr: "localhost:8080",
		},
		{
			name: "{host:\"192.168.1.1\", port:1234,} => 192.168.1.1:1234",
			addr: NetAddress{
				host: "localhost",
				port: 8080,
			},
			expectedStr: "localhost:8080",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedStr, test.addr.String())
		})
	}
}

func TestNetAddressSet(t *testing.T) {
	tests := []struct {
		name          string
		addr          NetAddress
		value         string
		expectedAddr  NetAddress
		isError       bool
		expectedError error
	}{
		{
			name:  "set localhost:8080",
			addr:  NetAddress{},
			value: "localhost:8080",
			expectedAddr: NetAddress{
				host: "localhost",
				port: 8080,
			},
			isError: false,
		},
		{
			name:          "error of set with no port",
			addr:          NetAddress{},
			value:         "localhost",
			isError:       true,
			expectedError: errors.New("need address in a form host:port"),
		},
		{
			name:          "error of set with bad port",
			addr:          NetAddress{},
			value:         "localhost:bad-port-number",
			isError:       true,
			expectedError: errors.New("port must be non negarive"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.isError {
				assert.Error(t, test.expectedError, test.addr.Set(test.value))
			} else {
				test.addr.Set(test.value)
				assert.Equal(t, test.expectedAddr, test.addr)
			}
		})
	}
}
