package agent

import (
	"github.com/k0st1a/metrics/internal/agent/client"
	"github.com/k0st1a/metrics/internal/metrics"
	"net/http"
)

func Run() {
	var myMetrics *metrics.MyStats = &metrics.MyStats{}
	var pollInternal = 2
	var reportInterval = 10
	var myClient *http.Client = &http.Client{}

	go metrics.RunUpdateMetrics(myMetrics, pollInternal)
	client.RunReportMetrics(myClient, myMetrics, reportInterval)
}
