package v18

import (
	"michelprogram/photon-parser/internal/context"
	"michelprogram/photon-parser/internal/reader"
)

// ReliableHeaderParameterCountV18 decodes reliable parameter counts for protocol v18.
type ReliableHeaderParameterCountV18 struct{}

var _ context.ReliableHeaderParameterCount = (*ReliableHeaderParameterCountV18)(nil)

// Count reads the reliable parameter count from reader.
func (ReliableHeaderParameterCountV18) Count(r *reader.Reader) (int, error) {
	res, err := r.ReadVarintUInt32()
	if err != nil {
		return 0, err
	}
	return int(res), nil
}
