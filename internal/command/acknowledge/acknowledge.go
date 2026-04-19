package acknowledge

import "michelprogram/photon-parser/internal/reader"

type Acknowledge struct{}

func (a Acknowledge) Parse(r *reader.Reader) error {
	return nil
}
