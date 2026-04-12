package photonparser

import "michelprogram/photon-parser/internal/reader"

type Reader = reader.Reader

func NewReader(data []byte) *Reader {
	return reader.NewReader(data)
}

type Parseable interface {
	Parse(r *Reader) error
}
