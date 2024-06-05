package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/rs/zerolog/log"
)

type Server struct {
	server *http.Server
}

func New(ctx context.Context, address string, handler http.Handler) *Server {
	s := &http.Server{
		BaseContext: func(_ net.Listener) context.Context { return ctx },
		Addr:        address,
		Handler:     handler,
	}

	return &Server{
		server: s,
	}
}

func (s *Server) Run() error {
	log.Printf("Run api")

	err := s.server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("listen and serve error:%w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	//nolint:wrapcheck //no need here
	return s.server.Shutdown(ctx)
}
