package agent

import (
	"github.com/k0st1a/metrics/internal/agent/report"
	"github.com/k0st1a/metrics/internal/metrics"
	"net/http"
)

func Run() {
	var myMetrics = &metrics.MyStats{}
	var myClient = &http.Client{}

	parseFlags()

	go metrics.RunUpdateMetrics(myMetrics, pollInternal)
	report.RunReportMetrics(serverAddr, myClient, myMetrics, reportInterval)
}
