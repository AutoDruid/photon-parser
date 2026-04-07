package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"golang.org/x/exp/constraints"
)

type Parseable interface {
	Parse(reader *bytes.Reader) error
}

type Parser[T any] func(reader *bytes.Reader) (*T, error)
type HeaderParser[H any] func(reader *bytes.Reader) (*H, error)

type Reader struct {
	*bytes.Reader
}

func NewReader(data []byte) *Reader {
	return &Reader{bytes.NewReader(data)}
}

func ReadHeader[H any](r *Reader) (*H, error) {
	var header H
	if err := binary.Read(r, binary.BigEndian, &header); err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}
	return &header, nil
}

func (r *Reader) ReadBytes(n int) ([]byte, error) {
	data := make([]byte, n)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, fmt.Errorf("failed to read %d bytes: %w", n, err)
	}
	return data, nil
}

func ReadPrimitive[T constraints.Integer | constraints.Float](r *Reader) (T, error) {
	var val T
	if err := binary.Read(r, binary.BigEndian, &val); err != nil {
		return val, fmt.Errorf("failed to read %T: %w", val, err)
	}
	return val, nil
}
