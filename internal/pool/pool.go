package pool

import (
	"sync"

	"michelprogram/photon-parser/internal/types"
)

// Params wraps a reusable parameter slice so sync.Pool can store
// a stable pointer and avoid per-round-trip heap allocation.
type Params struct {
	S []types.Parameter
}

var paramPool = sync.Pool{
	New: func() any {
		return &Params{S: make([]types.Parameter, 0, 16)}
	},
}

// Get returns a Params with S of length n, reusing the underlying array when possible.
func Get(n int) *Params {
	p := paramPool.Get().(*Params)
	if cap(p.S) < n {
		p.S = make([]types.Parameter, n)
	} else {
		p.S = p.S[:n]
		clear(p.S) // Go 1.21+
	}
	return p
}

// Put returns the Params to the pool. Caller must not reference p or p.S afterwards.
func Put(p *Params) {
	if p == nil {
		return
	}
	// Drop pathologically large slices so one weird packet doesn't bloat the pool.
	if cap(p.S) > 4096 {
		return
	}
	p.S = p.S[:0]
	paramPool.Put(p)
}
