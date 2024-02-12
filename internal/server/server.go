package server

import (
	"fmt"
	"net/http"

	"github.com/k0st1a/metrics/internal/handlers"
	"github.com/k0st1a/metrics/internal/handlers/json"
	"github.com/k0st1a/metrics/internal/handlers/text"
	"github.com/k0st1a/metrics/internal/storage/file"
	"github.com/k0st1a/metrics/internal/storage/inmemory"
	"github.com/rs/zerolog/log"
)

type Storage interface {
	GetGauge(string) (float64, bool)
	StoreGauge(string, float64)

	GetCounter(string) (int64, bool)
	StoreCounter(string, int64)

	GetAll() (map[string]int64, map[string]float64)
}

func Run() error {
	log.Debug().Msg("Run server")

	cfg, err := collectConfig()
	if err != nil {
		return err
	}

	printConfig(cfg)

	var s Storage
	if cfg.FileStoragePath == "" {
		log.Debug().Msg("Using memory storage")
		s = inmemory.NewStorage()
	} else {
		log.Debug().Msg("Using file storage")
		s = file.NewStorage(cfg.FileStoragePath, cfg.StoreInterval, cfg.Restore)
	}

	th := text.NewHandler(s)
	jh := json.NewHandler(s)

	r := handlers.NewRouter()
	text.BuildRouter(r, th)
	json.BuildRouter(r, jh)

	err = http.ListenAndServe(cfg.ServerAddr, r)
	if err != nil {
		return fmt.Errorf("listen and serve error:%w", err)
	}

	return nil
}
