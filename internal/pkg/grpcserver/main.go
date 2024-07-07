// Package grpcserver is some behaviour of GRPC server.
package grpcserver

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/rs/zerolog/log"
)

type Server struct {
	Server   *grpc.Server
	Listener *net.Listener
}

// New - создание сервера, где:
//   - address - хост и порт сервера;
//   - srv - запускаемое grpc api.
func New(address string, server *grpc.Server) (*Server, error) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("net listen error:%w", err)
	}

	return &Server{
		Listener: &l,
		Server:   server,
	}, nil
}

// Run - запуск сервера.
func (s *Server) Run() error {
	log.Printf("Run grpc api")

	err := s.Server.Serve(*s.Listener)
	if err != nil {
		return fmt.Errorf("server listen error:%w", err)
	}

	return nil
}

// Shutdown - graceful выключение сервера.
func (s *Server) Shutdown(_ context.Context) error {
	s.Server.GracefulStop()
	return nil
}
