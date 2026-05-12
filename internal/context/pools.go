package context

import (
	"sync"
)

const (
	maxPooledCap = 1024
)


type Pool[P any] struct {
	pool sync.Pool
}

func NewPool[P any](maxCap int) *Pool[P] {

	return &Pool[P]{
		pool: sync.Pool{
			New: func() any {
				return make([]P, 0, maxCap)
			},
		},
	}
}

func (p *Pool[P]) Get(n int) []P {
	if n < 0 {
		n = 0
	}
	raw := p.pool.Get()
	s, _ := raw.([]P)
	if cap(s) >= n {
		return s[:n]
	}
	return make([]P, n)
}

func (p *Pool[P]) Put(buff []P) {
	if cap(buff) > maxPooledCap {
		return
	}
	var zero P
	for i := range buff {
		buff[i] = zero
	}
	p.pool.Put(buff[:0])
}
