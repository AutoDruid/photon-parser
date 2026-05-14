package acknowledge

import (
	"encoding/binary"

	"github.com/AutoDruid/photon-parser/internal/reader"
	"github.com/AutoDruid/photon-parser/internal/types"
)

func Parse(reader *reader.Reader, out *types.Acknowledge) error {
	var err error

	out.AckReliableSequenceNumber, err = reader.ReadUInt32(binary.BigEndian)
	if err != nil {
		return err
	}

	out.AckSentTime, err = reader.ReadUInt32(binary.BigEndian)
	if err != nil {
		return err
	}

	return nil

}
