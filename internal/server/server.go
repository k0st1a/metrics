package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	hdbp "github.com/k0st1a/metrics/internal/handlers/db/ping"
	sdbs "github.com/k0st1a/metrics/internal/storage/db"
	sdbm "github.com/k0st1a/metrics/internal/storage/db/migration/v1"
	sdbp "github.com/k0st1a/metrics/internal/storage/db/ping"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/k0st1a/metrics/internal/handlers"
	"github.com/k0st1a/metrics/internal/handlers/json"
	"github.com/k0st1a/metrics/internal/handlers/text"
	"github.com/k0st1a/metrics/internal/middleware"
	"github.com/k0st1a/metrics/internal/middleware/checksign"
	"github.com/k0st1a/metrics/internal/middleware/decrypt"
	"github.com/k0st1a/metrics/internal/pkg/crypto/rsa"
	"github.com/k0st1a/metrics/internal/pkg/hash"
	"github.com/k0st1a/metrics/internal/pkg/profiler"
	"github.com/k0st1a/metrics/internal/pkg/retry"
	"github.com/k0st1a/metrics/internal/pkg/server"
	"github.com/k0st1a/metrics/internal/storage/file"
	"github.com/k0st1a/metrics/internal/storage/inmemory"
	"github.com/rs/zerolog/log"
)

type Storage interface {
	GetGauge(ctx context.Context, name string) (*float64, error)
	StoreGauge(ctx context.Context, name string, value float64) error

	GetCounter(ctx context.Context, name string) (*int64, error)
	StoreCounter(ctx context.Context, name string, value int64) error

	StoreAll(ctx context.Context, counter map[string]int64, gauge map[string]float64) error
	GetAll(ctx context.Context) (counter map[string]int64, gauge map[string]float64, err error)
}

type Pinger interface {
	Ping(ctx context.Context) error
}

func Run() error {
	log.Debug().Msg("Run server")

	cfg, err := NewConfig()
	if err != nil {
		return err
	}

	printConfig(cfg)

	var s Storage
	var p Pinger

	ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelCtx()

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

	rt := retry.New()
	th := text.NewHandler(s, rt)
	jh := json.NewHandler(s, rt)
	dbph := hdbp.NewHandler(p)

	var middlewares []func(http.Handler) http.Handler

	if cfg.HashKey != "" {
		h := hash.New(cfg.HashKey)
		middlewares = append(middlewares, checksign.New(h))
	}

	if cfg.CryptoKey != "" {
		prv, err := rsa.NewPrivateFromFile(cfg.CryptoKey)
		if err != nil {
			return fmt.Errorf("rsa new private from file error:%w", err)
		}

		middlewares = append(middlewares, decrypt.New(prv))
	}

	middlewares = append(middlewares, middleware.Logging, middleware.Compress)

	r := handlers.NewRouter(middlewares)

	text.BuildRouter(r, th)
	json.BuildRouter(r, jh)
	hdbp.BuildRouter(r, dbph)

	srv, err := server.New(ctx, cfg.ServerAddr, r)
	if err != nil {
		return fmt.Errorf("metrics server new error:%w", err)
	}

	go func() {
		err := srv.Run()
		if errors.Is(err, http.ErrServerClosed) {
			log.Printf("metrics server closed")
			return
		}
		if err != nil {
			log.Error().Err(err).Msg("failed to run metrics server")
		}
	}()

	prf, err := profiler.New(ctx, cfg.PprofServerAddr)
	if err != nil {
		return fmt.Errorf("profiler server new error:%w", err)
	}

	go func() {
		err := prf.Run()
		if errors.Is(err, http.ErrServerClosed) {
			log.Printf("profiler server closed")
			return
		}
		if err != nil {
			log.Error().Err(err).Msg("failed to run profiler server")
		}
	}()

	<-ctx.Done()

	err = srv.Shutdown(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("error of shutdown metrics server")
	}

	err = prf.Shutdown(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("error of shutdown profiler server")
	}

	return nil
}
