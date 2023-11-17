package server

import (
	"errors"
	"flag"
	"strconv"
	"strings"

	"github.com/k0st1a/metrics/internal/logger"
)

var flagRunAddr string

type NetAddress struct {
	host string
	port int
}

func parseFlags() error {
	addr := &NetAddress{host: "localhost", port: 8080}
	// если интерфейс не реализован,
	// здесь будет ошибка компиляции
	_ = flag.Value(addr)
	flag.Var(addr, "a", "server network address")
	flag.Parse()
	logger.Printf("Host:%s Port:%d\n", addr.host, addr.port)
	logger.Printf("Args:%v", flag.Args()[:])

	if len(flag.Args()) != 0 {
		return errors.New("unknown args")
	}

	flagRunAddr = addr.String()

	return nil
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
