// Package server is http server for process metrics from requests.
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	hping "github.com/k0st1a/metrics/internal/adapters/api/http/server/handler/db/ping"
	"github.com/k0st1a/metrics/internal/adapters/storage/db"
	v1 "github.com/k0st1a/metrics/internal/adapters/storage/db/migration/v1"
	dbping "github.com/k0st1a/metrics/internal/adapters/storage/db/ping"

	"github.com/3th1nk/cidr"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/k0st1a/metrics/internal/adapters/api/http/middleware/checksign"
	"github.com/k0st1a/metrics/internal/adapters/api/http/middleware/compress"
	"github.com/k0st1a/metrics/internal/adapters/api/http/middleware/decrypt"
	"github.com/k0st1a/metrics/internal/adapters/api/http/middleware/logging"
	"github.com/k0st1a/metrics/internal/adapters/api/http/middleware/trustedsubnet"
	"github.com/k0st1a/metrics/internal/adapters/api/http/server/config"
	"github.com/k0st1a/metrics/internal/adapters/api/http/server/handler"
	"github.com/k0st1a/metrics/internal/adapters/api/http/server/handler/json"
	"github.com/k0st1a/metrics/internal/adapters/api/http/server/handler/text"
	"github.com/k0st1a/metrics/internal/adapters/storage/file"
	"github.com/k0st1a/metrics/internal/adapters/storage/inmemory"
	"github.com/k0st1a/metrics/internal/pkg/crypto/rsa"
	"github.com/k0st1a/metrics/internal/pkg/hash"
	"github.com/k0st1a/metrics/internal/pkg/profiler"
	"github.com/k0st1a/metrics/internal/pkg/retry"
	"github.com/k0st1a/metrics/internal/pkg/server"
	"github.com/k0st1a/metrics/internal/ports"
	"github.com/rs/zerolog/log"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

func Run() error {
	log.Debug().Msg("Run server")

	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("config create error:%w", err)
	}

	log.Printf("Cfg:%+v", cfg)

	var s ports.Storage
	var p Pinger

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

	rt := retry.New()
	th := text.NewHandler(s, rt)
	jh := json.NewHandler(s, rt)
	dbph := hping.NewHandler(p)

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

	if cfg.TrustedSubnet != "" {
		subnet, err := cidr.Parse(cfg.TrustedSubnet)
		if err != nil {
			return fmt.Errorf("cidr parse trusted subnet error:%w", err)
		}

		middlewares = append(middlewares, trustedsubnet.New(subnet))
	}

	middlewares = append(middlewares, logging.New, compress.New)

	r := handler.NewRouter(middlewares)

	text.BuildRouter(r, th)
	json.BuildRouter(r, jh)
	hping.BuildRouter(r, dbph)

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
