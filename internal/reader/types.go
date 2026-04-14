package reader

type Parseable interface {
	Parse(r *Reader) error
}

type Payload interface{}

type Reader struct {
	Buffer []byte
	Max    int
	Cursor int
}

const (
	INT8_SIZE    = 1
	INT16_SIZE   = 2
	INT32_SIZE   = 4
	INT64_SIZE   = 8
	FLOAT32_SIZE = 4
	FLOAT64_SIZE = 8
)