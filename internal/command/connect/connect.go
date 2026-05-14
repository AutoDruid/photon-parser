package connect

import (
	"github.com/AutoDruid/photon-parser/internal/reader"
	"github.com/AutoDruid/photon-parser/internal/types"
)

func Parse(reader *reader.Reader, out *types.Connect) error {
	var err error

	out.Mtu, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}
	out.WindowSize, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}
	out.ChannelCount, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}
	out.IncomingBandwidth, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}
	out.OutgoingBandwidth, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}
	out.DisconnectThrottle, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}
	out.PacketThrottleAcceleration, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}
	out.PacketThrottleDeceleration, err = reader.ReadUInt32BE()
	if err != nil {
		return err
	}

	return nil
}
