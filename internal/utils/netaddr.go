package utils

import (
	"errors"
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
		return err
	}

	port16, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return errors.New("invalid port " + strconv.Quote(port) + " parsing " + strconv.Quote(flagValue)) // from netip.ParseAddrPort
	}

	a.host = host
	a.port = uint16(port16)

	return nil
}
