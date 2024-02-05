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

type reporter struct {
	d              Doer
	reportInterval int
}

func NewReporter(d Doer, i int) Runner {
	return &reporter{
		d:              d,
		reportInterval: i,
	}
}

func (r *reporter) Run() {
	ticker := time.NewTicker(time.Duration(r.reportInterval) * time.Second)

	for range ticker.C {
		r.d.Do()
	}
}
