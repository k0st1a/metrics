package server

import (
	"fmt"
	"net/http"

	"github.com/k0st1a/metrics/internal/handlers"
	"github.com/k0st1a/metrics/internal/handlers/json"
	"github.com/k0st1a/metrics/internal/handlers/text"
	"github.com/k0st1a/metrics/internal/storage"
)

func Run() error {
	cfg, err := collectConfig()
	if err != nil {
		return err
	}

	printConfig(cfg)

	r := handlers.NewRouter()

	s := storage.NewStorage()
	th := text.NewHandler(s)
	jh := json.NewHandler(s)

	text.BuildRouter(r, th)
	json.BuildRouter(r, jh)

	err = http.ListenAndServe(cfg.ServerAddr, r)
	if err != nil {
		return fmt.Errorf("listen and serve error:%w", err)
	}

	return nil
}
