package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type dbMigration struct {
	c *pgxpool.Pool
}

// NewMigration - создание сущности "миграция".
func NewMigration(c *pgxpool.Pool) *dbMigration {
	return &dbMigration{
		c: c,
	}
}

// Migrate - запускает миграцию.
func (db *dbMigration) Migrate(ctx context.Context) error {
	tx, err := db.c.Begin(ctx)
	if err != nil {
		return fmt.Errorf("db migration transaction begin error:%w", err)
	}
	defer func() {
		err = tx.Rollback(ctx)
		switch {
		case errors.Is(err, pgx.ErrTxClosed):
			log.Debug().Msg("db migration transaction closed")
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

	tag, err := tx.Exec(ctx, q1)
	if err != nil {
		return fmt.Errorf("db migration in transaction create counters error:%w", err)
	}
	log.Printf("tag of create counters table:%v", tag)

	tag, err = tx.Exec(ctx, `CREATE INDEX IF NOT EXISTS counter_idx ON counters (name)`)
	if err != nil {
		return fmt.Errorf("db migration in transaction create counters index error:%w", err)
	}
	log.Printf("tag of create counters index:%v", tag)

	q2 := `
        CREATE TABLE IF NOT EXISTS gauges(
                name  varchar(40)      PRIMARY KEY,
				value double precision NULL
        )
	`

	tag, err = tx.Exec(ctx, q2)
	if err != nil {
		return fmt.Errorf("db migration in transaction create gauges error:%w", err)
	}
	log.Printf("tag of create gauges table:%v", tag)

	tag, err = tx.Exec(ctx, `CREATE INDEX IF NOT EXISTS gauge_idx ON gauges (name)`)
	if err != nil {
		return fmt.Errorf("db migration in transaction create gauges index error:%w", err)
	}
	log.Printf("tag of create gauges index:%v", tag)

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("db migration transaction commit error:%w", err)
	}

	return nil
}
