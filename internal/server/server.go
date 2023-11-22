package server

import (
	"net/http"

	"github.com/k0st1a/metrics/internal/handlers"
	"github.com/k0st1a/metrics/internal/storage"
)

func Run() error {
	cfg := NewConfig()
	err := collectConfig(&cfg)
	if err != nil {
		return err
	}

	storage := storage.NewStorage()
	handler := handlers.NewHandler(storage)

	err = http.ListenAndServe(cfg.ServerAddr, handlers.BuildRouter(handler))
	if err != nil {
		return err
	}

	return nil
}
