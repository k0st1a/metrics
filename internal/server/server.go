package server

import (
	"net/http"

	"github.com/k0st1a/metrics/internal/handlers"
	"github.com/k0st1a/metrics/internal/logger"
	"github.com/k0st1a/metrics/internal/storage"
)

func Run() {
	cfg := NewConfig()
	logger.Println("Config:", cfg)

	err := parseFlags(&cfg)
	if err != nil {
		panic(err)
	}
	logger.Println("Config after parseFlags:", cfg)

	err = parseEnv(&cfg)
	if err != nil {
		panic(err)
	}
	logger.Println("Config after parseEnv:", cfg)

	logger.Println("Storage running")
	storage.Run()
	logger.Println("Storage runned")

	err = http.ListenAndServe(cfg.ServerAddr, handlers.BuildRouter())
	if err != nil {
		panic(err)
	}
}
