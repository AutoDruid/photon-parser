package v16

import "michelprogram/photon-parser/internal/reader"

type ReliableHeaderParameterCountV16 struct{}

var _ reader.ReliableHeaderParameterCount = (*ReliableHeaderParameterCountV16)(nil)

func (ReliableHeaderParameterCountV16) Count(r *reader.Reader) (int, error) {
	res, err := r.ReadUInt8()
	if err != nil{
		return 0, err
	}

	return int(res), nil
}
