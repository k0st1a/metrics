package ping

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Pinger interface {
	Ping(context.Context) error
}

type dbPing struct {
	c *pgxpool.Pool
}

func NewPinger(c *pgxpool.Pool) Pinger {
	return &dbPing{
		c: c,
	}
}

func (db *dbPing) Ping(ctx context.Context) error {
	err := db.c.Ping(ctx)
	if err != nil {
		return fmt.Errorf("db.c.ping error:%w", err)
	}

	return nil
}
