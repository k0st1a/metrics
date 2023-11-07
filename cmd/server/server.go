package main

import (
	"github.com/k0st1a/metrics/internal/logger"
	"github.com/k0st1a/metrics/internal/server"
)

func main() {
	logger.Run()
	server.Run()
}