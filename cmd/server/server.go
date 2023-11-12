package main

import (
	"fmt"

	"github.com/k0st1a/metrics/internal/logger"
	"github.com/k0st1a/metrics/internal/server"
)

func main() {
	fmt.Println("Running logger")
	logger.Run()
	defer logger.Close()

	fmt.Println("Running server")
	server.Run()
}