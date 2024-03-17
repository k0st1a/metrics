package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Migrator interface {
	Migrate(ctx context.Context) error
}

type dbMigration struct {
	c *pgxpool.Pool
}

func NewMigration(c *pgxpool.Pool) Migrator {
	return &dbMigration{
		c: c,
	}
}

func (db *dbMigration) Migrate(ctx context.Context) error {
	tx, err := db.c.Begin(ctx)
	if err != nil {
		return fmt.Errorf("db migration transaction begin error:%w", err)
	}
	defer func() {
		err = tx.Rollback(ctx)
		switch {
		case errors.Is(err, pgx.ErrTxClosed):
			log.Debug().Err(err).Msg("db migration transaction closed")
		default:
			log.Error().Err(err).Msg("db migration transaction close error")
		}
	}()

	q1 := `
        CREATE TABLE IF NOT EXISTS counters(
                name  varchar(40)      PRIMARY KEY,
				delta bigint           NULL
        )
	`

	t, err := tx.Exec(ctx, q1)
	if err != nil {
		return fmt.Errorf("db migration in transaction create counters error:%w", err)
	}
	log.Printf("t:%v", t)

	q2 := `
        CREATE TABLE IF NOT EXISTS gauges(
                name  varchar(40)      PRIMARY KEY,
				value double precision NULL
        )
	`

	t, err = tx.Exec(ctx, q2)
	if err != nil {
		return fmt.Errorf("db migration in transaction create gauges error:%w", err)
	}
	log.Printf("t:%v", t)

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("db migration transaction commit error:%w", err)
	}

	return nil
}
