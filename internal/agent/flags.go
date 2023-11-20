package agent

import (
	"errors"
	"flag"

	"github.com/k0st1a/metrics/internal/utils"
)

func parseFlags(cfg *Config) error {
	addr := &utils.NetAddress{}
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
