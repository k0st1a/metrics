package server

import (
	"net/http"

	"github.com/k0st1a/metrics/internal/logger"
	"github.com/k0st1a/metrics/internal/server/handlers"
)

func Run() {
	mux := http.NewServeMux()
	logger.Println("Mux started")

	mux.HandleFunc("/", handlers.Stub)
	logger.Println("Stab handler added")

	mux.HandleFunc("/update/gauge/", handlers.Gauge)
	logger.Println("Gauge handler added")

	mux.HandleFunc("/update/counter/", handlers.Counter)
	logger.Println("Counter handler added")

	err := http.ListenAndServe(`localhost:8080`, mux)
	if err != nil {
		panic(err)
	}
}
