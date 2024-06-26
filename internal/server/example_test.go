package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/k0st1a/metrics/internal/storage/db"
	v1 "github.com/k0st1a/metrics/internal/storage/db/migration/v1"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/k0st1a/metrics/internal/handlers"
	"github.com/k0st1a/metrics/internal/handlers/json"
	"github.com/k0st1a/metrics/internal/middleware"
	"github.com/k0st1a/metrics/internal/middleware/checksign"
	"github.com/k0st1a/metrics/internal/pkg/hash"
	"github.com/k0st1a/metrics/internal/pkg/retry"
	"github.com/k0st1a/metrics/internal/pkg/server"
	"github.com/rs/zerolog/log"
)

func Example() { //nolint:testableexamples // no output here
	cfg, _ := NewConfig()

	ctx := context.Background()

	pool, _ := pgxpool.New(ctx, cfg.DatabaseDSN)

	m := v1.NewMigration(pool)
	_ = m.Migrate(ctx)

	s := db.NewStorage(pool)

	rt := retry.New()
	jh := json.NewHandler(s, rt)

	var middlewares []func(http.Handler) http.Handler

	if cfg.HashKey != "" {
		h := hash.New(cfg.HashKey)
		middlewares = append(middlewares, checksign.New(h))
	}

	middlewares = append(middlewares, middleware.Logging, middleware.Compress)

	r := handlers.NewRouter(middlewares)

	json.BuildRouter(r, jh)

	srv, _ := server.New(ctx, cfg.ServerAddr, r)

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

	<-ctx.Done()

	err := srv.Shutdown(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("error of shutdown metrics server")
	}
}
