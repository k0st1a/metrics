// Package db for work with PostgreSQL DB.
package db

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/k0st1a/metrics/internal/utils"
	"github.com/rs/zerolog/log"
)

type DBStorage struct {
	c *pgxpool.Pool
	m sync.Mutex
}

// NewStorage - создать storage для хранения метрик в БД, где:
//   - c - пулл коннекций до БД.
func NewStorage(c *pgxpool.Pool) *DBStorage {
	return &DBStorage{
		c: c,
	}
}

// StoreGauge - сохраняет метрику типа gauge с именем name и значенем value.
func (s *DBStorage) StoreGauge(ctx context.Context, name string, value float64) error {
	s.m.Lock()
	defer s.m.Unlock()

	_, err := s.c.Exec(ctx, "INSERT INTO gauges (name,value) VALUES($1, $2)", name, value)
	if err != nil {
		return fmt.Errorf("store gauge query error:%w", err)
	}

	return nil
}

// GetGauge - возвращает метрику типа gauge с именем name.
func (s *DBStorage) GetGauge(ctx context.Context, name string) (*float64, error) {
	log.Printf("GetGauge, name:%v", name)
	var v float64

	err := s.c.QueryRow(ctx, "SELECT value FROM gauges WHERE name = $1 LIMIT 1", name).Scan(&v)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, utils.ErrMetricsNoGauge
	}
	if err != nil {
		return nil, fmt.Errorf("get gauge query error:%w", err)
	}

	return &v, nil
}

// StoreCounter - сохраняет метрику типа counter с именем name и значенем value.
func (s *DBStorage) StoreCounter(ctx context.Context, name string, value int64) error {
	s.m.Lock()
	defer s.m.Unlock()

	_, err := s.c.Exec(ctx, "INSERT INTO counters (name,delta) VALUES($1, $2)"+
		"ON CONFLICT (name) DO UPDATE SET delta = counters.delta + $2", name, value)
	if err != nil {
		return fmt.Errorf("store counter query error:%w", err)
	}

	return nil
}

// GetCounter - возвращает метрику типа gauge с именем name.
func (s *DBStorage) GetCounter(ctx context.Context, name string) (*int64, error) {
	log.Printf("GetCounter, name:%v", name)
	var d int64

	err := s.c.QueryRow(ctx, "SELECT delta FROM counters WHERE name = $1 LIMIT 1", name).Scan(&d)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, utils.ErrMetricsNoCounter
	}
	if err != nil {
		return nil, fmt.Errorf("get counter query error:%w", err)
	}

	return &d, nil
}

// StoreAll - сохраняет группу метрик типа counter и gauge.
func (s *DBStorage) StoreAll(ctx context.Context, counter map[string]int64, gauge map[string]float64) error {
	s.m.Lock()
	defer s.m.Unlock()

	var b pgx.Batch

	log.Printf("StoreAll, counter:%v gauge:%v", counter, gauge)

	for k, v := range counter {
		b.Queue("INSERT INTO counters (name,delta) VALUES($1, $2)"+
			"ON CONFLICT (name) DO UPDATE SET delta = counters.delta + $2", k, v)
	}

	for k2, v2 := range gauge {
		b.Queue("INSERT INTO gauges (name,value) VALUES($1, $2) ON CONFLICT (name) DO UPDATE SET value = $2", k2, v2)
	}

	br := s.c.SendBatch(ctx, &b)
	defer func() {
		err := br.Close()
		if err != nil {
			log.Error().Err(err).Msg("store all br close error")
		}
	}()

	for i := 1; i < len(counter)+len(gauge); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("br.exec error:%w", err)
		}
	}

	return nil
}

// GetAll - возвращает все метрики типа counter и gauge.
func (s *DBStorage) GetAll(ctx context.Context) (map[string]int64, map[string]float64, error) {
	var b pgx.Batch

	b.Queue("SELECT name,delta FROM counters")
	b.Queue("SELECT name,value FROM gauges")

	br := s.c.SendBatch(ctx, &b)
	defer func() {
		err := br.Close()
		if err != nil {
			log.Error().Err(err).Msg("get all br close error")
		}
	}()

	rows, err := br.Query()
	if err != nil {
		return nil, nil, fmt.Errorf("get counters query error:%w", err)
	}
	defer rows.Close()

	c := make(map[string]int64)

	for rows.Next() {
		var name string
		var delta int64

		err = rows.Scan(&name, &delta)
		if err != nil {
			return nil, nil, fmt.Errorf("counter rows scan error:%w", err)
		}

		c[name] = delta
	}

	err = rows.Err()
	if err != nil {
		return nil, nil, fmt.Errorf("counter rows error:%w", err)
	}

	rows, err = br.Query()
	if err != nil {
		return nil, nil, fmt.Errorf("get gauges query error:%w", err)
	}
	defer rows.Close()

	g := make(map[string]float64)

	for rows.Next() {
		var name string
		var value float64

		err = rows.Scan(&name, &value)
		if err != nil {
			return nil, nil, fmt.Errorf("gauge rows scan error:%w", err)
		}

		g[name] = value
	}

	err = rows.Err()
	if err != nil {
		return nil, nil, fmt.Errorf("gauge rows error:%w", err)
	}

	return c, g, nil
}
