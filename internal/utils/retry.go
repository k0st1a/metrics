package utils

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type Retryer interface {
	Retry(ctx context.Context, check func(error) bool, fnc func() error) error
}

type retry struct {
	intervals []time.Duration
}

func NewRetry() Retryer {
	return &retry{
		intervals: []time.Duration{
			time.Duration(1) * time.Second,
			//nolint //2-ой повтор через 3 секунды
			time.Duration(3) * time.Second,
			//nolint //3-ой повтор через 5 секунд
			time.Duration(5) * time.Second,
		},
	}
}

func (r *retry) Retry(ctx context.Context, check func(error) bool, fnc func() error) error {
	err := fnc()

	for _, interval := range r.intervals {
		if check(err) {
			err = wait(ctx, interval)
			if err != nil {
				return err
			}
			err = fnc()
			continue
		}

		return err
	}

	return ErrMaxRetryReached
}

func IsConnectionException(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code)
}

func wait(ctx context.Context, interval time.Duration) error {
	timer := time.NewTimer(interval)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		//nolint //Возвращаем ошибку завершкуния контекста
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
