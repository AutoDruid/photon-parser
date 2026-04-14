package ping

import "michelprogram/photon-parser/internal/reader"

type Ping struct{}

var _ reader.Parseable = (*Ping)(nil)

func (p Ping) Parse(r *reader.Reader) error {
	return nil
}