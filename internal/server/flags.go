package server

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
	flag.Var(addr, "a", "server network address")
	flag.Parse()

	if len(flag.Args()) != 0 {
		return errors.New("unknown args")
	}

	cfg.ServerAddr = addr.String()

	return nil
}
