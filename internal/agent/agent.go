package agent

import (
	"github.com/k0st1a/metrics/internal/agent/report"
	"github.com/k0st1a/metrics/internal/metrics"
	"net/http"
)

func Run() {
	var myMetrics = &metrics.MyStats{}
	var myClient = &http.Client{}
	var pollInternal = 2
	var reportInterval = 10

	go metrics.RunUpdateMetrics(myMetrics, pollInternal)
	report.RunReportMetrics(myClient, myMetrics, reportInterval)
}
