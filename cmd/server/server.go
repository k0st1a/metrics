package main

import (
	"github.com/k0st1a/metrics/internal/server"
	"github.com/rs/zerolog/log"
)

func main() {
	//log.Logger = log.With().Caller().Logger()
	//name := fmt.Sprintf("./myapp-%v.log", time.Now().String())
	//log.Printf("Log to file:%v", name)
	//file, _ := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	//log.Logger = zerolog.New(file).With().Caller().Timestamp().Logger()

	err := server.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
