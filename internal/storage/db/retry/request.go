package retry

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/k0st1a/metrics/internal/utils"
)

type Storage interface {
	StoreGauge(ctx context.Context, name string, value float64) error
	GetGauge(ctx context.Context, name string) (*float64, error)

	StoreCounter(ctx context.Context, name string, value int64) error
	GetCounter(ctx context.Context, name string) (*int64, error)

	StoreAll(ctx context.Context, counter map[string]int64, gauge map[string]float64) error
	GetAll(ctx context.Context) (counter map[string]int64, gauge map[string]float64, err error)
}

type dbRetry struct {
	storage        Storage
	retryIntervals []time.Duration
}

func NewRetry(s Storage) Storage {
	return &dbRetry{
		storage: s,
		retryIntervals: []time.Duration{
			time.Duration(1) * time.Second,
			//nolint //2-ой повтор через 3 секунды
			time.Duration(3) * time.Second,
			//nolint //3-ой повтор через 5 секунд
			time.Duration(5) * time.Second,
		},
	}
}

func (r *dbRetry) StoreGauge(ctx context.Context, name string, value float64) (err error) {
	err = r.storage.StoreGauge(ctx, name, value)

	for _, interval := range r.retryIntervals {
		if isConnectionException(err) {
			time.Sleep(interval)
			err = r.storage.StoreGauge(ctx, name, value)
			continue
		}

		//nolint // Возвращаем значение как есть
		return err
	}

	return utils.ErrMaxRetryReached
}

func (r *dbRetry) GetGauge(ctx context.Context, name string) (value *float64, err error) {
	value, err = r.storage.GetGauge(ctx, name)

	for _, interval := range r.retryIntervals {
		if isConnectionException(err) {
			time.Sleep(interval)
			value, err = r.storage.GetGauge(ctx, name)
			continue
		}

		//nolint // Возвращаем значения как есть
		return value, err
	}

	return nil, utils.ErrMaxRetryReached
}

func (r *dbRetry) StoreCounter(ctx context.Context, name string, value int64) (err error) {
	err = r.storage.StoreCounter(ctx, name, value)

	for _, interval := range r.retryIntervals {
		if isConnectionException(err) {
			time.Sleep(interval)
			err = r.storage.StoreCounter(ctx, name, value)
			continue
		}

		//nolint // Возвращаем значение как есть
		return err
	}

	return utils.ErrMaxRetryReached
}
func (r *dbRetry) GetCounter(ctx context.Context, name string) (value *int64, err error) {
	value, err = r.storage.GetCounter(ctx, name)

	for _, interval := range r.retryIntervals {
		if isConnectionException(err) {
			time.Sleep(interval)
			value, err = r.storage.GetCounter(ctx, name)
			continue
		}

		//nolint // Возвращаем значения как есть
		return value, err
	}

	return nil, utils.ErrMaxRetryReached
}

func (r *dbRetry) StoreAll(ctx context.Context, counter map[string]int64, gauge map[string]float64) (err error) {
	err = r.storage.StoreAll(ctx, counter, gauge)

	for _, interval := range r.retryIntervals {
		if isConnectionException(err) {
			time.Sleep(interval)
			err = r.storage.StoreAll(ctx, counter, gauge)
			continue
		}

		//nolint // Возвращаем значение как есть
		return err
	}

	return utils.ErrMaxRetryReached
}

func (r *dbRetry) GetAll(ctx context.Context) (counter map[string]int64, gauge map[string]float64, err error) {
	counter, gauge, err = r.storage.GetAll(ctx)

	for _, interval := range r.retryIntervals {
		if isConnectionException(err) {
			time.Sleep(interval)
			counter, gauge, err = r.storage.GetAll(ctx)
			continue
		}

		//nolint // Возвращаем значения как есть
		return counter, gauge, err
	}

	return nil, nil, utils.ErrMaxRetryReached
}

func isConnectionException(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code)
}
