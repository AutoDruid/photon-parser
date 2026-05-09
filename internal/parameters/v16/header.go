package v16

import (
	"michelprogram/photon-parser/internal/context"
	"michelprogram/photon-parser/internal/reader"
)

// ReliableHeaderParameterCountV16 decodes reliable parameter counts for protocol v16.
type ReliableHeaderParameterCountV16 struct{}

var _ context.ReliableHeaderParameterCount = (*ReliableHeaderParameterCountV16)(nil)

// Count reads the reliable parameter count from reader.
func (ReliableHeaderParameterCountV16) Count(r *reader.Reader) (int, error) {
	res, err := r.ReadUInt8()
	if err != nil {
		return 0, err
	}

	return int(res), nil
}
