package pool

import (
	"michelprogram/photon-parser/internal/types"
	"sync"
)

type Params[P types.VersionedParameter] struct {
	S []P
}

type Pool[P types.VersionedParameter] struct {
	pool sync.Pool
}

func New[P types.VersionedParameter]() *Pool[P] {
	return &Pool[P]{
		pool: sync.Pool{
			New: func() any {
				return &Params[P]{S: make([]P, 0, 16)}
			},
		},
	}
}

func (p *Pool[P]) Get(n int) *Params[P] {
	params := p.pool.Get().(*Params[P])
	if cap(params.S) < n {
		params.S = make([]P, n)
	} else {
		params.S = params.S[:n]
		clear(params.S)
	}
	return params
}

func (p *Pool[P]) Put(params *Params[P]) {
	if params == nil {
		return
	}
	if cap(params.S) > 4096 {
		return
	}
	params.S = params.S[:0]
	p.pool.Put(params)
}
