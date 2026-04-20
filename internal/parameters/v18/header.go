package v18

import "michelprogram/photon-parser/internal/reader"

type ReliableHeaderParameterCountV18 struct{}

var _ reader.ReliableHeaderParameterCount = (*ReliableHeaderParameterCountV18)(nil)

func (ReliableHeaderParameterCountV18) Count(r *reader.Reader) (int, error) {
	res, err := r.ReadVarintInt32()
	if err != nil{
		return 0, err
	}

	return int(res), nil
}
