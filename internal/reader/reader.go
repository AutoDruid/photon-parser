package reader

import (
	"github.com/AutoDruid/photon-parser/internal/errors"
)

const (
	INT8_SIZE       = 1
	INT16_SIZE      = 2
	INT32_SIZE      = 4
	INT64_SIZE      = 8
	FLOAT32_SIZE    = 4
	FLOAT64_SIZE    = 8
	VARINT_SHIFT    = 7
	VARINT_MASK     = 0x7F
	VARINT_MSB_MASK = 0x80
)

type Reader struct {
	Buffer []byte
	Max    int
	Cursor int
}

func NewReader(data []byte) *Reader {
	return &Reader{
		Buffer: data,
		Max:    len(data),
		Cursor: 0,
	}
}

// ReadRemaining reads the remaining bytes from the reader.
// This is useful when you want to read the rest of the data in the reader.
func (r *Reader) ReadRemaining() []byte {
	tmp := r.Cursor
	r.Cursor = r.Max
	return r.Buffer[tmp:]
}

// Skip n bytes from the reader.
// Doing this is faster than reading the bytes and discarding them.
func (r *Reader) Skip(n int) error {

	if n < 0 {
		return errors.ErrInvalidNegativeSkip
	}

	size := r.Cursor + n

	if size > r.Max {
		return errors.ErrNotEnoughBytesBytes
	}

	r.Cursor += n
	return nil
}

// Reset resets the reader to the beginning of the buffer.
func (r *Reader) Reset(data []byte) {
	r.Buffer = data
	r.Cursor = 0
	r.Max = len(data)
}
