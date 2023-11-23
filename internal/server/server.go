package server

import (
	"net/http"

	"github.com/k0st1a/metrics/internal/handlers"
	"github.com/k0st1a/metrics/internal/storage"
)

func Run() error {
	cfg, err := collectConfig()
	if err != nil {
		return err
	}

	printConfig(cfg)

	storage := storage.NewStorage()
	handler := handlers.NewHandler(storage)

	err = http.ListenAndServe(cfg.ServerAddr, handlers.BuildRouter(handler))
	if err != nil {
		return err
	}

	return nil
}
