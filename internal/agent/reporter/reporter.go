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
	doer           Doer
	reportInterval int
}

func NewReporter(d Doer, i int) Runner {
	return &reporter{
		doer:           d,
		reportInterval: i,
	}
}

func (r *reporter) Run() {
	ticker := time.NewTicker(time.Duration(r.reportInterval) * time.Second)

	for range ticker.C {
		r.doer.Do()
	}
}
