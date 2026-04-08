// Package parser provides core binary parsing utilities for the Photon Protocol.
// It includes type-safe generic readers and helper types for reading binary data
// in big-endian format.
package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"sync"
)

// Number is a type constraint for all numeric types supported by the binary reader.
type Number interface {
	~int8 | ~int16 | ~int32 | ~int64 |
		~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Reader wraps a bytes.Reader to provide additional parsing utilities.
// All data is read in big-endian format as per Photon Protocol specification.
type Reader struct {
	*bytes.Reader
}

// readerPool is a pool of Reader instances to reduce allocations on hot paths.
var readerPool = sync.Pool{
	New: func() any {
		return &Reader{new(bytes.Reader)}
	},
}

// NewReader creates a new Reader from a byte slice.
// The underlying bytes.Reader will read data in sequential order.
// Prefer NewReaderFromPool / ReleaseReader on performance-sensitive paths.
func NewReader(data []byte) *Reader {
	return &Reader{bytes.NewReader(data)}
}

// NewReaderFromPool retrieves a Reader from the shared pool and resets it to
// read from data. The caller must call ReleaseReader when done.
func NewReaderFromPool(data []byte) *Reader {
	r := readerPool.Get().(*Reader)
	r.Reset(data)
	return r
}

// ReleaseReader returns a pooled Reader back to the pool.
// The Reader must have been obtained with NewReaderFromPool.
func ReleaseReader(r *Reader) {
	readerPool.Put(r)
}

// ReadHeader reads a struct header of type H from the reader using binary.Read.
// The header type H must be a struct with exported fields that can be read
// via encoding/binary in big-endian format.
//
// Example:
//
//	type MyHeader struct {
//	    Version uint16
//	    Length  uint32
//	}
//	header, err := ReadHeader[MyHeader](reader)
func ReadHeader[H any](r *Reader) (*H, error) {
	var header H
	if err := binary.Read(r, binary.BigEndian, &header); err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}
	return &header, nil
}

// ReadBytes reads exactly n bytes from the reader.
// It returns an error if fewer than n bytes are available.
// This uses io.ReadFull internally to ensure all bytes are read.
func (r *Reader) ReadBytes(n int) ([]byte, error) {
	data := make([]byte, n)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, fmt.Errorf("failed to read %d bytes: %w", n, err)
	}
	return data, nil
}

// ReadPrimitive reads a primitive type T from the reader in big-endian format.
// T must be a numeric type (integer or float) that encoding/binary can handle.
// This is a generic function that works with int8, int16, int32, int64,
// uint8, uint16, uint32, uint64, float32, and float64.
//
// Example:
//
//	value, err := ReadPrimitive[int32](reader)
func ReadPrimitive[T Number](r *Reader) (T, error) {
	var val T
	if err := binary.Read(r, binary.BigEndian, &val); err != nil {
		return val, fmt.Errorf("failed to read %T: %w", val, err)
	}
	return val, nil
}
