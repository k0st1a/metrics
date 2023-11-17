package server

import (
	"net/http"

	"github.com/k0st1a/metrics/internal/handlers"
	"github.com/k0st1a/metrics/internal/logger"
	"github.com/k0st1a/metrics/internal/storage"
)

func Run() {
	parseFlags()

	logger.Println("Storage running")
	storage.Run()
	logger.Println("Storage runned")

	err := http.ListenAndServe(flagRunAddr, handlers.BuildRouter())
	if err != nil {
		panic(err)
	}
}
