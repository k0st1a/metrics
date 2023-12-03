package server

import (
	"fmt"
	"net/http"

	"github.com/k0st1a/metrics/internal/handlers"
	"github.com/k0st1a/metrics/internal/handlers/json"
	"github.com/k0st1a/metrics/internal/handlers/text"
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

	r := handlers.NewRouter()

	s := storage.NewStorage()
	gs := gauge.NewGaugeStorage(s)
	cs := counter.NewCounterStorage(s)

	csh := text.NewHandler(cs)
	gsh := text.NewHandler(gs)

	mh := json.NewHandler(s)

	text.BuildRouter(r, csh, gsh)
	json.BuildRouter(r, mh)

	err = http.ListenAndServe(cfg.ServerAddr, r)
	if err != nil {
		return fmt.Errorf("listen and serve error:%w", err)
	}

	return nil
}
