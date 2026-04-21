package connect

import (
	"michelprogram/photon-parser/internal/hooks"
	"michelprogram/photon-parser/internal/reader"
)

type Connect struct {
	Mtu                        uint32
	WindowSize                 uint32
	ChannelCount               uint32
	IncomingBandwidth          uint32
	OutgoingBandwidth          uint32
	DisconnectThrollte         uint32
	PacketThrottleAcceleration uint32
	PacketThrottleDeceleration uint32
}

func Parse(r *reader.Reader, hooks *hooks.Hooks) (*Connect, error) {
	var err error
	connect := Connect{}
	connect.Mtu, err = r.ReadUInt32()
	if err != nil {
		return nil, err
	}
	connect.WindowSize, err = r.ReadUInt32()
	if err != nil {
		return nil, err
	}
	connect.ChannelCount, err = r.ReadUInt32()
	if err != nil {
		return nil, err
	}
	connect.IncomingBandwidth, err = r.ReadUInt32()
	if err != nil {
		return nil, err
	}
	connect.OutgoingBandwidth, err = r.ReadUInt32()
	if err != nil {
		return nil, err
	}
	connect.DisconnectThrollte, err = r.ReadUInt32()
	if err != nil {
		return nil, err
	}
	connect.PacketThrottleAcceleration, err = r.ReadUInt32()
	if err != nil {
		return nil, err
	}
	connect.PacketThrottleDeceleration, err = r.ReadUInt32()
	if err != nil {
		return nil, err
	}
	return &connect, nil
}
