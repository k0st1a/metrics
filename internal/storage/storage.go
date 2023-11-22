package storage

import "github.com/rs/zerolog/log"

var gauge map[string]float64
var counter map[string]int64

func RunStorage() {
	log.Debug().Msg("RunStorage")

	gauge = make(map[string]float64)
	counter = make(map[string]int64)
}

func StoreGauge(name string, value float64) {
	log.Debug().
		Str("name:", name).
		Float64("value", value).
		Msg("StoreGauge")

	gauge[name] = value
}

func GetGauge(name string) (float64, bool) {
	v, ok := gauge[name]
	log.Debug().
		Str("name:", name).
		Float64("value:", v).
		Msg("GetGauge")

	return v, ok
}

func StoreCounter(name string, value int64) {
	log.Debug().
		Str("name", name).
		Int64("value", value).
		Msg("StoreCounter")

	counter[name] = value
}

func GetCounter(name string) (int64, bool) {
	v, ok := counter[name]
	log.Debug().
		Str("name", name).
		Int64("value", v).
		Msg("GetCounter")

	return v, ok
}
