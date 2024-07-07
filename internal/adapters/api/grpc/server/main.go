// Package server for process request by grpc.
package server

import (
	"fmt"

	"google.golang.org/grpc"

	pb "github.com/k0st1a/metrics/internal/adapters/api/grpc/protobuf"
	"github.com/k0st1a/metrics/internal/adapters/api/grpc/server/handler"
	"github.com/k0st1a/metrics/internal/application/server/config"
	"github.com/k0st1a/metrics/internal/pkg/grpcserver"
	"github.com/k0st1a/metrics/internal/pkg/retry"
	"github.com/k0st1a/metrics/internal/ports"
)

func New(cfg *config.Config, storage ports.Storage) (*grpcserver.Server, error) {
	rt := retry.New()

	h := &handler.MetricsServer{
		Storage: storage,
		Retry:   rt,
	}

	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer()

	// регистрируем сервис
	pb.RegisterMetricsServer(s, h)

	srv, err := grpcserver.New(cfg.ServerAddr, s)
	if err != nil {
		return nil, fmt.Errorf("grpc metrics server new error:%w", err)
	}

	return srv, nil
}
