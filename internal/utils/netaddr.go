package utils

import (
	"fmt"
	"net"
	"strconv"
)

type NetAddress struct {
	host string
	port uint16
}

func (a *NetAddress) String() string {
	return a.host + ":" + strconv.FormatUint(uint64(a.port), 10)
}

func (a *NetAddress) Set(flagValue string) error {
	host, port, err := net.SplitHostPort(flagValue)
	if err != nil {
		return fmt.Errorf("host:port split error:%w", err)
	}

	port16, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return fmt.Errorf("port parsing error:%w", err)
	}

	a.host = host
	a.port = uint16(port16)

	return nil
}
