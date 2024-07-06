// Package server - пакет для получения/хранения метрик.
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/k0st1a/metrics/internal/adapters/storage/db"
	v1 "github.com/k0st1a/metrics/internal/adapters/storage/db/migration/v1"
	dbping "github.com/k0st1a/metrics/internal/adapters/storage/db/ping"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/k0st1a/metrics/internal/adapters/api/http/server"
	"github.com/k0st1a/metrics/internal/adapters/storage/file"
	"github.com/k0st1a/metrics/internal/adapters/storage/inmemory"
	"github.com/k0st1a/metrics/internal/application/server/config"
	"github.com/k0st1a/metrics/internal/pkg/profiler"
	"github.com/k0st1a/metrics/internal/ports"
	"github.com/rs/zerolog/log"
)

func Run() error {
	log.Debug().Msg("Run server")

	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("config create error:%w", err)
	}

	log.Printf("Cfg:%+v", cfg)

	var s ports.Storage
	var p ports.Pinger

	ctx, cancelFunc := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancelFunc()

	switch {
	case cfg.DatabaseDSN != "":
		log.Debug().Msg("Using db storage")
		pool, err := pgxpool.New(ctx, cfg.DatabaseDSN)
		if err != nil {
			return fmt.Errorf("pgxpool new error:%w", err)
		}

		m := v1.NewMigration(pool)
		err = m.Migrate(ctx)
		if err != nil {
			return fmt.Errorf("migrate error:%w", err)
		}

		p = dbping.NewPinger(pool)
		s = db.NewStorage(pool)

	case cfg.FileStoragePath != "":
		log.Debug().Msg("Using file storage")
		s = file.NewStorage(ctx, cfg.FileStoragePath, cfg.StoreInterval, cfg.Restore)

	default:
		log.Debug().Msg("Using memory storage")
		s = inmemory.NewStorage()
	}

	srv, err := server.New(ctx, cfg, s, p)
	if err != nil {
		return fmt.Errorf("make server error:%w", err)
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
