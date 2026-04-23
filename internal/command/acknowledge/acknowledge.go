package acknowledge

import (
	"encoding/binary"
	"michelprogram/photon-parser/internal/context"
)

type Acknowledge struct {
	AckReliableSequenceNumber uint32
	AckSentTime               uint32
}

func Parse(ctx *context.Context) (*Acknowledge, error) {
	var ack Acknowledge
	var err error

	ack.AckReliableSequenceNumber, err = ctx.Reader.ReadUInt32(binary.BigEndian)
	if err != nil {
		return nil, err
	}

	ack.AckSentTime, err = ctx.Reader.ReadUInt32(binary.BigEndian)
	if err != nil {
		return nil, err
	}

	return &ack, nil

}
