package report

import (
	"github.com/k0st1a/metrics/internal/metrics"
	"io"
	"net/http"
	"os"
	"time"
)

func RunReportMetrics(addr string, client *http.Client, metrics *metrics.MyStats, reportInterval int) {
	for {
		time.Sleep(time.Duration(reportInterval) * time.Second)
		reportMetrics(addr, client, metrics)
	}
}

func reportMetrics(addr string, client *http.Client, metrics *metrics.MyStats) {
	prepared := metrics.Prepare()
	for name, info := range prepared {
		ReportMetric(addr, client, info.Type, name, info.Value)
	}
}

func ReportMetric(addr string, client *http.Client, metricType, name, value string) {
	var url = `http://` + addr + `/update/` + metricType + `/` + name + `/` + value
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Content-Length", "0")

	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, response.Body) // вывод ответа в консоль
	response.Body.Close()
}
