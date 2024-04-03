package reporter

import (
	"net/http"
	"time"

	"github.com/k0st1a/metrics/internal/agent/model"
	"github.com/k0st1a/metrics/internal/agent/report/json"
)

type state struct {
	pollerCh       chan<- struct{}
	serverAddr     string
	reportInterval int
	rateLimit      int
	sign           http.RoundTripper
}

//nolint:lll //no need here
func NewReporter(serverAddr string, reportInterval int, rateLimit int, sign http.RoundTripper) (*state, <-chan struct{}) {
	pollerCh := make(chan struct{})
	return &state{
		serverAddr:     serverAddr,
		reportInterval: reportInterval,
		rateLimit:      rateLimit,
		pollerCh:       pollerCh,
	}, pollerCh
}

func (s *state) Do(reportCh <-chan map[string]model.MetricInfoRaw) {
	agentCh := make(chan map[string]model.MetricInfoRaw)
	for i := 0; i < s.rateLimit; i++ {
		c := &http.Client{
			Transport: s.sign,
		}
		go json.NewReport(s.serverAddr, c, agentCh).Do()
	}

	reportTicker := time.NewTicker(time.Duration(s.reportInterval) * time.Second)

	for range reportTicker.C {
		s.pollerCh <- struct{}{}
		m := <-reportCh
		if len(m) == 0 {
			continue
		}

		agentCh <- m
	}
}
