// Package server is some behaviour of HTTP server.
package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/rs/zerolog/log"
)

type Server struct {
	Server   *http.Server
	Listener *net.Listener
}

// New - создание сервера, где:
//   - ctx - контекст отмены запросов, обрабатываемых сервером;
//   - address - хост и порт сервера;
//   - handler - обработчик сервера.
func New(ctx context.Context, address string, handler http.Handler) (*Server, error) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("net listen error:%w", err)
	}

	s := &http.Server{
		BaseContext: func(_ net.Listener) context.Context { return ctx },
		Addr:        l.Addr().String(),
		Handler:     handler,
	}

	return &Server{
		Server:   s,
		Listener: &l,
	}, nil
}

// Run - запуск сервера.
func (s *Server) Run() error {
	log.Printf("Run api")

	err := s.Server.Serve(*s.Listener)
	if err != nil {
		return fmt.Errorf("server listen error:%w", err)
	}

	return nil
}

// Shutdown - graceful выключение сервера.
func (s *Server) Shutdown(ctx context.Context) error {
	//nolint:wrapcheck //no need here
	return s.Server.Shutdown(ctx)
}
