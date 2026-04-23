package connect

import (
	"michelprogram/photon-parser/internal/context"
)

type Connect struct {
	Mtu                        uint32
	WindowSize                 uint32
	ChannelCount               uint32
	IncomingBandwidth          uint32
	OutgoingBandwidth          uint32
	DisconnectThrottle         uint32
	PacketThrottleAcceleration uint32
	PacketThrottleDeceleration uint32
}

func Parse(ctx *context.Context) (*Connect, error) {
	var err error
	connect := Connect{}
	connect.Mtu, err = ctx.Reader.ReadUInt32()
	if err != nil {
		return nil, err
	}
	connect.WindowSize, err = ctx.Reader.ReadUInt32()
	if err != nil {
		return nil, err
	}
	connect.ChannelCount, err = ctx.Reader.ReadUInt32()
	if err != nil {
		return nil, err
	}
	connect.IncomingBandwidth, err = ctx.Reader.ReadUInt32()
	if err != nil {
		return nil, err
	}
	connect.OutgoingBandwidth, err = ctx.Reader.ReadUInt32()
	if err != nil {
		return nil, err
	}
	connect.DisconnectThrottle, err = ctx.Reader.ReadUInt32()
	if err != nil {
		return nil, err
	}
	connect.PacketThrottleAcceleration, err = ctx.Reader.ReadUInt32()
	if err != nil {
		return nil, err
	}
	connect.PacketThrottleDeceleration, err = ctx.Reader.ReadUInt32()
	if err != nil {
		return nil, err
	}
	return &connect, nil
}
