package acknowledge

import (
	"encoding/binary"

	"github.com/AutoDruid/photon-parser/internal/reader"
)

type Acknowledge struct {
	AckReliableSequenceNumber uint32
	AckSentTime               uint32
}

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
