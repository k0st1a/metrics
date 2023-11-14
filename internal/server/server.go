package server

import (
	"net/http"

	"github.com/k0st1a/metrics/internal/handlers"
	"github.com/k0st1a/metrics/internal/logger"
	"github.com/k0st1a/metrics/internal/storage"
)

func Run() {
	logger.Println("Run storage")
	storage.Run()

	err := http.ListenAndServe(":8080", handlers.BuildRouter())
	if err != nil {
		panic(err)
	}
}
