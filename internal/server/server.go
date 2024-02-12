package server

import (
	"fmt"
	"net/http"

	"github.com/k0st1a/metrics/internal/handlers"
	dbh "github.com/k0st1a/metrics/internal/handlers/db"
	"github.com/k0st1a/metrics/internal/handlers/json"
	"github.com/k0st1a/metrics/internal/handlers/text"
	dbs "github.com/k0st1a/metrics/internal/storage/db"
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

type Pinger interface {
	Ping() error
}

func Run() error {
	log.Debug().Msg("Run server")

	cfg, err := collectConfig()
	if err != nil {
		return err
	}

	printConfig(cfg)

	var p Pinger
	if cfg.DatabaseDSN != "" {
		p = dbs.NewStorage(cfg.DatabaseDSN)
	}

	var s Storage
	switch {
	case cfg.FileStoragePath != "":
		log.Debug().Msg("Using file storage")
		s = file.NewStorage(cfg.FileStoragePath, cfg.StoreInterval, cfg.Restore)
	default:
		log.Debug().Msg("Using memory storage")
		s = inmemory.NewStorage()
	}

	th := text.NewHandler(s)
	jh := json.NewHandler(s)
	dh := dbh.NewHandler(p)

	r := handlers.NewRouter()
	text.BuildRouter(r, th)
	json.BuildRouter(r, jh)
	dbh.BuildRouter(r, dh)

	err = http.ListenAndServe(cfg.ServerAddr, r)
	if err != nil {
		return fmt.Errorf("listen and serve error:%w", err)
	}

	return nil
}
