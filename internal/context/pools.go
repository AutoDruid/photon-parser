package context

import (
	"sync"

	"github.com/AutoDruid/photon-parser/internal/types"
)

const (
	maxPooledCap = 1024
)

type Pool[P types.ParameterView] struct {
	pool sync.Pool
}

type PooledSlice[P types.ParameterView] struct{
	Items []P
}

func NewPool[P types.ParameterView](maxCap int) *Pool[P] {
	return &Pool[P]{
		pool: sync.Pool{
			New: func() any {
				return &PooledSlice[P]{
					Items: make([]P, 0, maxCap),
				}
			},
		},
	}
}

func (p *Pool[P]) Get(n int) *PooledSlice[P] {
	if n < 0 {
		n = 0
	}
	wrapper := p.pool.Get().(*PooledSlice[P])
	
	if cap(wrapper.Items) >= n {
		wrapper.Items = wrapper.Items[:n]
		return wrapper
	}
	wrapper.Items = make([]P, n)
	return wrapper
}

func (p *Pool[P]) Put(wrapper *PooledSlice[P]) {
	if cap(wrapper.Items) > maxPooledCap {
		return
	}
	wrapper.Items = wrapper.Items[:0]
	p.pool.Put(wrapper)
}
