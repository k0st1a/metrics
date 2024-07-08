// Package ratelimit to limit the number of agents.
package ratelimit

import "sync"

type state struct {
	limit int
}

func New(l int) *state {
	return &state{
		limit: l,
	}
}

func (s *state) Run(fn func()) {
	var wg sync.WaitGroup

	for i := 0; i < s.limit; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fn()
		}()
	}

	wg.Wait()
}
