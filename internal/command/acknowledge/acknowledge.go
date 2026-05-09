package acknowledge

import (
	"encoding/binary"
	"michelprogram/photon-parser/internal/reader"
)

// Acknowledge represents the payload of an acknowledge command.
type Acknowledge struct {
	AckReliableSequenceNumber uint32
	AckSentTime               uint32
}

// Parse decodes an acknowledge payload from reader.
func Parse(reader *reader.Reader) (*Acknowledge, error) {
	var ack Acknowledge
	var err error

	ack.AckReliableSequenceNumber, err = reader.ReadUInt32(binary.BigEndian)
	if err != nil {
		return nil, err
	}

	ack.AckSentTime, err = reader.ReadUInt32(binary.BigEndian)
	if err != nil {
		return nil, err
	}

	return &ack, nil

}
