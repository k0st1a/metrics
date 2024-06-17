// Package profiler for run pprof profiler.
package profiler

import (
	"context"
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

	return server.New(ctx, address, mux)
}
