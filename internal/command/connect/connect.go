package connect

import (
	"github.com/AutoDruid/photon-parser/internal/reader"
	"github.com/AutoDruid/photon-parser/internal/types"
)

func ParseInto(reader *reader.Reader, dest *types.Connect) error {
	var err error

	dest.Mtu, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}
	dest.WindowSize, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}
	dest.ChannelCount, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}
	dest.IncomingBandwidth, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}
	dest.OutgoingBandwidth, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}
	dest.DisconnectThrottle, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}
	dest.PacketThrottleAcceleration, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}
	dest.PacketThrottleDeceleration, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}

	return nil
}
