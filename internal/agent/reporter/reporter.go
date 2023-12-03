package reporter

import (
	"net/http"
	"time"

	"github.com/k0st1a/metrics/internal/metrics"
)

type Reporter interface {
	DoReportsMetrics(*http.Client, *metrics.MyStats)
}

type reporter struct {
	rtype Reporter
}

func NewReporter(r Reporter) reporter {
	return reporter{
		rtype: r,
	}
}

func (r *reporter) RunReportMetrics(c *http.Client, m *metrics.MyStats, ri int) {
	ticker := time.NewTicker(time.Duration(ri) * time.Second)

	for range ticker.C {
		r.reportMetrics(c, m)
	}
}

func (r *reporter) reportMetrics(c *http.Client, m *metrics.MyStats) {
	m.IncreasePollCount()
	r.rtype.DoReportsMetrics(c, m)
}
