package io

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

type IntervalWriter interface {
	Run(ctx context.Context, interval int)
}

type intervalWriter struct {
	writer  Writer
	storage StorageGeter
}

func NewIntervalWriter(w Writer, s StorageGeter) IntervalWriter {
	log.Debug().Msg("NewIntervalWriter")
	return &intervalWriter{
		writer:  w,
		storage: s,
	}
}

func (w *intervalWriter) Run(ctx context.Context, interval int) {
	log.Debug().Msg("Run interval writer")
	ticker := time.NewTicker(time.Duration(interval) * time.Second)

	for range ticker.C {
		log.Debug().Msg("Tick of interval writer")
		err := w.writer.Write(ctx, w.storage)
		if err != nil {
			log.Error().Err(err).Msg("write error storage to file")
		}
		log.Debug().Msg("Storage writed to file")
	}
}
