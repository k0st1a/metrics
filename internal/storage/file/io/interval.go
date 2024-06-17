// Package io для создания писателя, который
// периодически сохраняет метрики на файловую систему.
package io

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

type intervalWriter struct {
	writer  Writer
	storage StorageGeter
}

// NewIntervalWriter - создание писателя, который с заданной периодичностью сохраняет текущие метрики на
// файловую систему, где:
//   - w - интерфейс записи метрик;
//   - s - интерфейс получения метрик.
func NewIntervalWriter(w Writer, s StorageGeter) *intervalWriter {
	log.Debug().Msg("NewIntervalWriter")
	return &intervalWriter{
		writer:  w,
		storage: s,
	}
}

// Run - запуск писателя, где:
//   - ctx - контекст отмены записи;
//   - interval - интервал времени в секундах, через которые сохранять метрики на файловую систему.
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
