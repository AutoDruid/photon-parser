package v16

import (
	"AutoDruid/photon-parser/internal/context"
	"AutoDruid/photon-parser/internal/reader"
)

type ReliableHeaderParameterCountV16 struct{}

var _ context.ReliableHeaderParameterCount = (*ReliableHeaderParameterCountV16)(nil)

func (ReliableHeaderParameterCountV16) Count(r *reader.Reader) (int, error) {
	res, err := r.ReadUInt8()
	if err != nil {
		return 0, err
	}

	return int(res), nil
}
