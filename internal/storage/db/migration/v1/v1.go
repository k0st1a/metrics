package v1

import (
	"context"
	"fmt"
	//"time"

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
	//time.Sleep(1 * time.Second) // Чтобы база успела подняться. На сколько это правильно?

	q1 := `
        CREATE TABLE IF NOT EXISTS counters(
                name  varchar(40)      PRIMARY KEY,
				delta bigint           NULL
        )
	`

	t, err := db.c.Exec(ctx, q1)
	if err != nil {
		return fmt.Errorf("create counters error:%w", err)
	}
	log.Printf("t:%v", t)

	q2 := `
        CREATE TABLE IF NOT EXISTS gauges(
                name  varchar(40)      PRIMARY KEY,
				value double precision NULL
        )
	`

	t, err = db.c.Exec(ctx, q2)
	if err != nil {
		return fmt.Errorf("create gauges error:%w", err)
	}
	log.Printf("t:%v", t)

	return nil
}
