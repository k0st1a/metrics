package utils

import (
	"errors"
	"strconv"
	"strings"
)

type NetAddress struct {
	host string
	port int
}

func (a *NetAddress) String() string {
	return a.host + ":" + strconv.Itoa(a.port)
}

func (a *NetAddress) Set(flagValue string) error {
	pl := strings.Split(flagValue, ":")
	if len(pl) != 2 {
		return errors.New("need address in a form host:port")
	}

	port, err := strconv.Atoi(pl[1])
	switch {
	case err != nil:
		return err
	case port < 0:
		return errors.New("port must be non negarive")
	}

	a.host = pl[0]
	a.port = port
	return nil
}
