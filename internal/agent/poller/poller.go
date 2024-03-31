package poller

import "time"

type Updater interface {
	Update()
}

type Runner interface {
	Run()
}

type poller struct {
	updater      Updater
	pollInterval int
}

func NewPoller(u Updater, i int) Runner {
	return &poller{
		updater:      u,
		pollInterval: i,
	}
}

func (p *poller) Run() {
	ticker := time.NewTicker(time.Duration(p.pollInterval) * time.Second)

	for range ticker.C {
		p.updater.Update()
	}
}
