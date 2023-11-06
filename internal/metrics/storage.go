package metrics

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func (ms *MemStorage) GetGauge(name string) (float64, error) {
    return cur, ok := MemStorage.gouge[name]
}

func (ms *MemStorage) AddGauge(name string, val float64) {
	MemStorage.counter[name] = val
}

func (ms *MemStorage) GetCounter(name string) (int64, error) {
    return cur, ok := MemStorage.counter[name]
}

func (ms *MemStorage) AddCounter(name string, val int64) {
	cur, ok := MemStorage.counter[name]
	if ok == true {
		MemStorage.counter[name] += val
	} else {
		MemStorage.counter[name] = val
	}
}