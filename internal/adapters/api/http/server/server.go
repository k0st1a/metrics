// Package server for create http server.
package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/3th1nk/cidr"
	"github.com/k0st1a/metrics/internal/adapters/api/http/middleware/checksign"
	"github.com/k0st1a/metrics/internal/adapters/api/http/middleware/compress"
	"github.com/k0st1a/metrics/internal/adapters/api/http/middleware/decrypt"
	"github.com/k0st1a/metrics/internal/adapters/api/http/middleware/logging"
	"github.com/k0st1a/metrics/internal/adapters/api/http/middleware/trustedsubnet"
	"github.com/k0st1a/metrics/internal/adapters/api/http/server/handler"
	hping "github.com/k0st1a/metrics/internal/adapters/api/http/server/handler/db/ping"
	"github.com/k0st1a/metrics/internal/adapters/api/http/server/handler/json"
	"github.com/k0st1a/metrics/internal/adapters/api/http/server/handler/text"
	"github.com/k0st1a/metrics/internal/application/server/config"
	"github.com/k0st1a/metrics/internal/pkg/crypto/rsa"
	"github.com/k0st1a/metrics/internal/pkg/hash"
	"github.com/k0st1a/metrics/internal/pkg/retry"
	"github.com/k0st1a/metrics/internal/pkg/server"
	"github.com/k0st1a/metrics/internal/ports"
)

func New(ctx context.Context, cfg *config.Config, s ports.Storage, p ports.Pinger) (*server.Server, error) {
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
			return nil, fmt.Errorf("rsa new private from file error:%w", err)
		}

		middlewares = append(middlewares, decrypt.New(prv))
	}

	if cfg.TrustedSubnet != "" {
		subnet, err := cidr.Parse(cfg.TrustedSubnet)
		if err != nil {
			return nil, fmt.Errorf("cidr parse trusted subnet error:%w", err)
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
		return nil, fmt.Errorf("metrics server new error:%w", err)
	}

	return srv, nil
}
