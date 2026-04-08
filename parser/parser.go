// Package parser provides core binary parsing utilities for the Photon Protocol.
// It includes type-safe generic readers and helper types for reading binary data
// in big-endian format.
package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"golang.org/x/exp/constraints"
)

// Parseable is an interface for types that can parse themselves from a byte reader.
type Parseable interface {
	Parse(reader *bytes.Reader) error
}

// Parser is a generic function type that parses a type T from a byte reader.
type Parser[T any] func(reader *bytes.Reader) (*T, error)

// HeaderParser is a generic function type that parses a header of type H from a byte reader.
type HeaderParser[H any] func(reader *bytes.Reader) (*H, error)

// Reader wraps a bytes.Reader to provide additional parsing utilities.
// All data is read in big-endian format as per Photon Protocol specification.
type Reader struct {
	*bytes.Reader
}

// NewReader creates a new Reader from a byte slice.
// The underlying bytes.Reader will read data in sequential order.
func NewReader(data []byte) *Reader {
	return &Reader{bytes.NewReader(data)}
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
func ReadPrimitive[T constraints.Integer | constraints.Float](r *Reader) (T, error) {
	var val T
	if err := binary.Read(r, binary.BigEndian, &val); err != nil {
		return val, fmt.Errorf("failed to read %T: %w", val, err)
	}
	return val, nil
}
