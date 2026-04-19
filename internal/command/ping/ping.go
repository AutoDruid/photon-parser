package ping

import "michelprogram/photon-parser/internal/reader"

type Ping struct{}

func (p Ping) Parse(r *reader.Reader) error {
	return nil
}