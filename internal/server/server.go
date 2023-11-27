package server

import (
	"fmt"
	"net/http"

	"github.com/k0st1a/metrics/internal/handlers"
	"github.com/k0st1a/metrics/internal/storage"
	"github.com/k0st1a/metrics/internal/storage/counter"
	"github.com/k0st1a/metrics/internal/storage/gauge"
)

func Run() error {
	cfg, err := collectConfig()
	if err != nil {
		return err
	}

	printConfig(cfg)

	s := storage.NewStorage()
	gs := gauge.NewGaugeStorage(s)
	cs := counter.NewCounterStorage(s)
	csh := handlers.NewHandler(cs)
	gsh := handlers.NewHandler(gs)

	err = http.ListenAndServe(cfg.ServerAddr, handlers.BuildRouter(csh, gsh))
	if err != nil {
		return fmt.Errorf("listen and serve error:%w", err)
	}

	return nil
}
