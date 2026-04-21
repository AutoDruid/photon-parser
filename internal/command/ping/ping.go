package ping

import "michelprogram/photon-parser/internal/reader"

type Ping struct{}

func Parse(r *reader.Reader) (*Ping, error) {
	return &Ping{}, nil
}
