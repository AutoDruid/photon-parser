package acknowledge

import "michelprogram/photon-parser/internal/reader"

type Acknowledge struct{}

var _ reader.Parseable = (*Acknowledge)(nil)

func (a Acknowledge) Parse(r *reader.Reader) error {
	return nil
}
