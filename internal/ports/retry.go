package ports

import "context"

// Retryer - интерфейс повторного выполнения запроса в случае, если функция check возвращает true.
type Retryer interface {
	Retry(ctx context.Context, check func(error) bool, fnc func() error) error
}
