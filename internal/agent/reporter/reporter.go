package reporter

import (
	"time"
)

type Runner interface {
	Run()
}

type Doer interface {
	Do()
}

type IncreasePollCounter interface {
	IncreasePollCount()
}

type reporter struct {
	d              Doer
	reportInterval int
	m              IncreasePollCounter
}

func NewReporter(d Doer, i int, m IncreasePollCounter) Runner {
	return &reporter{
		d:              d,
		reportInterval: i,
		m:              m,
	}
}

func (r *reporter) Run() {
	ticker := time.NewTicker(time.Duration(r.reportInterval) * time.Second)

	for range ticker.C {
		r.m.IncreasePollCount()
		r.d.Do()
	}
}
