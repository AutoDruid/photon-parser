package acknowledge

import (
	"github.com/AutoDruid/photon-parser/internal/reader"
	"github.com/AutoDruid/photon-parser/internal/types"
)

func ParseInto(reader *reader.Reader, dest *types.Acknowledge) error {
	var err error

	dest.AckReliableSequenceNumber, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}

	dest.AckSentTime, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}

	return nil

}
