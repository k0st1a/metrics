package agent

import (
	"errors"
	"flag"
	"strconv"
	"strings"
)

type NetAddress struct {
	host string
	port int
}

func parseFlags(cfg *Config) error {
	addr := &NetAddress{}
	addr.Set(cfg.ServerAddr)

	// если интерфейс не реализован,
	// здесь будет ошибка компиляции
	_ = flag.Value(addr)
	flag.Var(addr, "a", "server network address in a form host:port")

	flag.IntVar(&(cfg.PollInterval), "p", 2, "metrics polling rate to the server")
	flag.IntVar(&(cfg.ReportInterval), "r", 10, "frequency of sending metrics to the server")

	flag.Parse()
	cfg.ServerAddr = addr.String()

	if len(flag.Args()) != 0 {
		return errors.New("unknown args")
	}

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
