package main

import (
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"math/rand"
)

type MetricForReport struct {
	Type  string
	Value string
}

type MyStats struct {
	PollCount   uint64
	RandomValue float64
	MemStats    runtime.MemStats
}

func (metrics *MyStats) update() {
	runtime.ReadMemStats(&metrics.MemStats)
	metrics.RandomValue = rand.Float64()
}

func (metrics *MyStats) Report(client *http.Client) {
	metrics.PollCount += 1

	preparedMetrics := map[string]MetricForReport{
		"Alloc":         MetricForReport{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.Alloc, 10)},
		"BuckHashSys":   MetricForReport{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.BuckHashSys, 10)},
		"Frees":         MetricForReport{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.Frees, 10)},
		"GCSys":         MetricForReport{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.GCSys, 10)},
		"HeapAlloc":     MetricForReport{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.HeapAlloc, 10)},
		"HeapIdle":      MetricForReport{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.HeapIdle, 10)},
		"HeapInuse":     MetricForReport{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.HeapInuse, 10)},
		"HeapObjects":   MetricForReport{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.HeapObjects, 10)},
		"HeapReleased":  MetricForReport{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.HeapReleased, 10)},
		"HeapSys":       MetricForReport{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.HeapSys, 10)},
		"LastGC":        MetricForReport{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.LastGC, 10)},
		"Lookups":       MetricForReport{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.Lookups, 10)},
		"MCacheInuse":   MetricForReport{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.MCacheInuse, 10)},
		"MCacheSys":     MetricForReport{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.MCacheSys, 10)},
		"MSpanInuse":    MetricForReport{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.MSpanInuse, 10)},
		"MSpanSys":      MetricForReport{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.MSpanSys, 10)},
		"Mallocs":       MetricForReport{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.Mallocs, 10)},
		"NextGC":        MetricForReport{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.NextGC, 10)},
		"OtherSys":      MetricForReport{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.OtherSys, 10)},
		"PauseTotalNs":  MetricForReport{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.PauseTotalNs, 10)},
		"StackInuse":    MetricForReport{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.StackInuse, 10)},
		"StackSys":      MetricForReport{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.StackSys, 10)},
		"Sys":           MetricForReport{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.Sys, 10)},
		"TotalAlloc":    MetricForReport{Type: "gauge", Value: strconv.FormatUint(metrics.MemStats.TotalAlloc, 10)},
		"NumForcedGC":   MetricForReport{Type: "gauge", Value: strconv.FormatUint(uint64(metrics.MemStats.NumForcedGC), 10)},
		"NumGC":         MetricForReport{Type: "gauge", Value: strconv.FormatUint(uint64(metrics.MemStats.NumGC), 10)},
		"GCCPUFraction": MetricForReport{Type: "gauge", Value: strconv.FormatFloat(metrics.MemStats.GCCPUFraction, 'g', -1, 64)},

		"PollCount":   MetricForReport{Type: "counter", Value: strconv.FormatUint(metrics.PollCount, 10)},
		"RandomValue": MetricForReport{Type: "gauge", Value: strconv.FormatFloat(metrics.RandomValue, 'g', -1, 64)},
	}

	for name, info := range preparedMetrics {
		ReportMetric(client, info.Type, name, info.Value)
	}
}

func RunUpdateMetrics(metrics *MyStats, pollInternal int) {
	for {
		metrics.update()
		time.Sleep(time.Duration(pollInternal) * time.Second)
	}
}

func RunReportMetrics(client *http.Client, metrics *MyStats, reportInterval int) {
	for {
		time.Sleep(time.Duration(reportInterval) * time.Second)
		metrics.Report(client)
	}
}

func ReportMetric(client *http.Client, metricType, name, value string) {
	var url string = `http://localhost:8080/update/` + metricType + `/` + name + `/` + value
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

func main() {
	var metrics *MyStats = &MyStats{}
	var pollInternal = 2
	var reportInterval = 10
	var client *http.Client = &http.Client{}

	go RunUpdateMetrics(metrics, pollInternal)
	RunReportMetrics(client, metrics, reportInterval)
}
