package storage

import "github.com/k0st1a/metrics/internal/logger"

var gauge map[string]float64
var counter map[string]int64

func Run() {
	gauge = make(map[string]float64)
	counter = make(map[string]int64)
}

func StoreGauge(name string, value float64) {
	logger.Println("StoreGauge, name:", name, "value:", value)
	gauge[name] = value
}

func GetGauge(name string) (float64, bool) {
	v, ok := gauge[name]
	logger.Println("GetGauge, name:", name, "value:", v)
	return v, ok
}

func StoreCounter(name string, value int64) {
	logger.Println("StoreCounter, name:", name, "value:", value)
	counter[name] = value
}

func GetCounter(name string) (int64, bool) {
	v, ok := counter[name]
	logger.Println("GetCounter, name:", name, "value:", v)
	return v, ok
}
