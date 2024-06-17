// Package profiler for run pprof profiler.
package profiler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/k0st1a/metrics/internal/pkg/server"
)

// New - создание сервера и обработчиками pprof.
func New(ctx context.Context, address string) (*server.Server, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	srv, err := server.New(ctx, address, mux)
	if err != nil {
		return nil, fmt.Errorf("profile server new error:%w", err)
	}

	return srv, nil
}
