package retry

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrMaxRetryReached = errors.New("retry: maximum number of retry reached")

type retry struct {
	intervals []time.Duration
}

// New - создание ретрайера, повтореное выполнение функции в зависимости от возвращаемой ею ошибки.
func New() *retry {
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

// Retry - запуск ретрайера, где:
// * ctx - контекст для отмены выполнения ретрайера;
// * check - функция проверки ошибки выполнения функции fnc;
// * fnc - данная фукнция выполняется повторно, если функция check возвращает true.
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

// IsConnectionException - проверка ошибки соединения с БД PostgreSQL
func IsConnectionException(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code)
}

func wait(ctx context.Context, interval time.Duration) error {
	timer := time.NewTimer(interval)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		//nolint //Возвращаем ошибку завершения контекста
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
