package v18

import (
	"AutoDruid/photon-parser/internal/context"
	"AutoDruid/photon-parser/internal/reader"
)

type ReliableHeaderParameterCountV18 struct{}

var _ context.ReliableHeaderParameterCount = (*ReliableHeaderParameterCountV18)(nil)

func (ReliableHeaderParameterCountV18) Count(r *reader.Reader) (int, error) {
	res, err := r.ReadVarintUInt32()
	if err != nil {
		return 0, err
	}
	return int(res), nil
}
