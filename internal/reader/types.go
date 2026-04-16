package reader

import "michelprogram/photon-parser/internal/types"

type Parseable interface {
	Parse(r *Reader) error
}
type Reader struct {
	Buffer []byte
	Max    int
	Cursor int

	types.SyncHooks
	types.AsyncHooks
}

const (
	INT8_SIZE    = 1
	INT16_SIZE   = 2
	INT32_SIZE   = 4
	INT64_SIZE   = 8
	FLOAT32_SIZE = 4
	FLOAT64_SIZE = 8
)
