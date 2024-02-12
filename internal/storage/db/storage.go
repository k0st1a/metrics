package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type Pinger interface {
	Ping() error
}

type dbStorage struct {
	databaseDSN string
}

func NewStorage(dsn string) Pinger {
	return &dbStorage{
		databaseDSN: dsn,
	}
}

func (db *dbStorage) Ping() error {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, db.databaseDSN)
	if err != nil {
		log.Error().Err(err).Msg("pgx.Connect error")
		return fmt.Errorf("pgx.Connect error:%w", err)
	}
	defer func() {
		err := conn.Close(ctx)
		if err != nil {
			log.Error().Err(err).Msg("conn.Close error")
		}
	}()

	err = conn.Ping(ctx)
	if err != nil {
		log.Error().Err(err).Msg("pgx.Ping error")
		return fmt.Errorf("pgx.Ping error:%w", err)
	}

	return nil
}
