// Package ports contains Storage interface to storage.
package ports

import "context"

// Storage - интерфейс работы с хранилищем метрик.
type Storage interface {
	// GetGauge - возвращает метрику типа gauge с именем name.
	GetGauge(ctx context.Context, name string) (*float64, error)
	// StoreGauge - сохраняет метрику типа gauge с именем name и значенем value.
	StoreGauge(ctx context.Context, name string, value float64) error

	// GetCounter - возвращает метрику типа gauge с именем name.
	GetCounter(ctx context.Context, name string) (*int64, error)
	// StoreCounter - сохраняет метрику типа counter с именем name и значенем value.
	StoreCounter(ctx context.Context, name string, value int64) error

	// StoreAll - сохраняет группу метрик типа counter и gauge.
	StoreAll(ctx context.Context, counter map[string]int64, gauge map[string]float64) error
	// GetAll - возвращает все метрики типа counter и gauge.
	GetAll(ctx context.Context) (counter map[string]int64, gauge map[string]float64, err error)
}
