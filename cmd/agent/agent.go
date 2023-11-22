package main

import (
	"github.com/k0st1a/metrics/internal/agent"
	"github.com/rs/zerolog/log"
)

func main() {
	err := agent.Run()
	if err != nil {
		log.Error().Err(err).Msg("")
	}
}
