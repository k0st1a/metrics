package utils

import (
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
		name         string
		addr         NetAddress
		value        string
		expectedAddr NetAddress
		isError      bool
		errString    string
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
			name:      "error of set with no port",
			addr:      NetAddress{},
			value:     "localhost",
			isError:   true,
			errString: "host:port split error:address localhost: missing port in address",
		},
		{
			name:      "error of set with bad port",
			addr:      NetAddress{},
			value:     "localhost:bad-port-number",
			isError:   true,
			errString: "port parsing error:strconv.ParseUint: parsing \"bad-port-number\": invalid syntax",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.isError {
				assert.EqualError(t, test.addr.Set(test.value), test.errString)
			} else {
				test.addr.Set(test.value)
				assert.Equal(t, test.expectedAddr, test.addr)
			}
		})
	}
}
