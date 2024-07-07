// Package client for create http client.
package client

import (
	"fmt"
	"net/http"

	"github.com/k0st1a/metrics/internal/adapters/api/http/client/json"
	"github.com/k0st1a/metrics/internal/adapters/api/http/middleware/encrypt"
	"github.com/k0st1a/metrics/internal/adapters/api/http/middleware/realip"
	"github.com/k0st1a/metrics/internal/adapters/api/http/middleware/roundtrip"
	"github.com/k0st1a/metrics/internal/adapters/api/http/middleware/sign"
	"github.com/k0st1a/metrics/internal/application/agent/config"
	"github.com/k0st1a/metrics/internal/pkg/crypto/rsa"
	"github.com/k0st1a/metrics/internal/pkg/hash"
	"github.com/k0st1a/metrics/internal/pkg/routing"
)

func New(cfg *config.Config) (*json.Client, error) {
	var middlewares []roundtrip.Middleware

	if cfg.HashKey != "" {
		h := hash.New(cfg.HashKey)
		middlewares = append(middlewares, sign.New(h))
	}

	if cfg.CryptoKey != "" {
		pbl, err := rsa.NewPublicFromFile(cfg.CryptoKey)
		if err != nil {
			return nil, fmt.Errorf("rsa new public from file error:%w", err)
		}

		middlewares = append(middlewares, encrypt.New(pbl))
	}

	host, err := routing.ParseHost(cfg.ServerAddr)
	if err != nil {
		return nil, fmt.Errorf("parse server address error:%w", err)
	}

	router, err := routing.New()
	if err != nil {
		return nil, fmt.Errorf("create router error:%w", err)
	}

	src, err := router.Route(host)
	if err != nil {
		return nil, fmt.Errorf("route to server address error:%w", err)
	}

	middlewares = append(middlewares, realip.New(src.String()))

	rt := roundtrip.New(http.DefaultTransport, middlewares...)

	c := json.New(cfg.ServerAddr, rt)

	return c, nil
}
