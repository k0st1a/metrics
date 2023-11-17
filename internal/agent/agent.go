package agent

import (
	"github.com/k0st1a/metrics/internal/agent/report"
	//"github.com/k0st1a/metrics/internal/logger"
	"fmt"
	"github.com/k0st1a/metrics/internal/metrics"
	"net/http"
)

func Run() {
	var myMetrics = &metrics.MyStats{}
	var myClient = &http.Client{}

	cfg := NewConfig()
	fmt.Println("Config:", cfg)
	parseFlags(&cfg)
	fmt.Println("Config after parseFlags:", cfg)
	parseEnv(&cfg)
	fmt.Println("Config after parseEnv:", cfg)

	go metrics.RunUpdateMetrics(myMetrics, cfg.PollInterval)
	report.RunReportMetrics(cfg.ServerAddr, myClient, myMetrics, cfg.ReportInterval)
}
