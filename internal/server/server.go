package server

import (
	"context"
	"fmt"
	"net/http"

	hdbp "github.com/k0st1a/metrics/internal/handlers/db/ping"
	sdbs "github.com/k0st1a/metrics/internal/storage/db"
	sdbm "github.com/k0st1a/metrics/internal/storage/db/migration/v1"
	sdbp "github.com/k0st1a/metrics/internal/storage/db/ping"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/k0st1a/metrics/internal/handlers"
	"github.com/k0st1a/metrics/internal/handlers/json"
	"github.com/k0st1a/metrics/internal/handlers/text"
	"github.com/k0st1a/metrics/internal/storage/file"
	"github.com/k0st1a/metrics/internal/storage/inmemory"
	"github.com/rs/zerolog/log"
)

type Storage interface {
	GetGauge(ctx context.Context, name string) (*float64, error)
	StoreGauge(ctx context.Context, name string, value float64) error

	GetCounter(ctx context.Context, name string) (*int64, error)
	StoreCounter(ctx context.Context, name string, value int64) error

	GetAll(ctx context.Context) (counter map[string]int64, gauge map[string]float64, err error)
}

type Pinger interface {
	Ping(ctx context.Context) error
}

func Run() error {
	log.Debug().Msg("Run server")

	cfg, err := collectConfig()
	if err != nil {
		return err
	}

	printConfig(cfg)

	var s Storage
	var p Pinger

	ctx := context.Background()

	switch {
	case cfg.DatabaseDSN != "":
		log.Debug().Msg("Using db storage")
		pool, err := pgxpool.New(ctx, cfg.DatabaseDSN)
		if err != nil {
			return fmt.Errorf("pgxpool new error:%w", err)
		}

		m := sdbm.NewMigration(pool)
		err = m.Migrate(ctx)
		if err != nil {
			return fmt.Errorf("migrate error:%w", err)
		}

		p = sdbp.NewPinger(pool)
		s = sdbs.NewStorage(pool)

	case cfg.FileStoragePath != "":
		log.Debug().Msg("Using file storage")
		s = file.NewStorage(ctx, cfg.FileStoragePath, cfg.StoreInterval, cfg.Restore)

	default:
		log.Debug().Msg("Using memory storage")
		s = inmemory.NewStorage()
	}

	th := text.NewHandler(s)
	jh := json.NewHandler(s)
	dbph := hdbp.NewHandler(p)

	r := handlers.NewRouter()
	text.BuildRouter(r, th)
	json.BuildRouter(r, jh)
	hdbp.BuildRouter(r, dbph)

	err = http.ListenAndServe(cfg.ServerAddr, r)
	if err != nil {
		return fmt.Errorf("listen and serve error:%w", err)
	}

	return nil
}
