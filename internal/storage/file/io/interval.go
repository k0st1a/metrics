package io

import (
	"time"

	"github.com/rs/zerolog/log"
)

type IntervalWriter interface {
	Run(i int)
}

type intervalWriter struct {
	writer  Writer
	storage Storage
}

func NewIntervalWriter(w Writer, s Storage) IntervalWriter {
	log.Debug().Msg("NewIntervalWriter")
	return &intervalWriter{
		writer:  w,
		storage: s,
	}
}

func (w *intervalWriter) Run(i int) {
	log.Debug().Msg("Run interval writer")
	ticker := time.NewTicker(time.Duration(i) * time.Second)

	for range ticker.C {
		log.Debug().Msg("Tick of interval writer")
		err := w.writer.Write(w.storage)
		if err != nil {
			log.Error().Err(err).Msg("write error storage to file")
		}
		log.Debug().Msg("Storage writed to file")
	}
}
