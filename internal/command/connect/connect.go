package connect

import (
	"encoding/binary"
	"michelprogram/photon-parser/internal/reader"
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

func Parse(reader *reader.Reader) (*Connect, error) {
	var err error
	connect := Connect{}
	connect.Mtu, err = reader.ReadUInt32(binary.BigEndian)
	if err != nil {
		return nil, err
	}
	connect.WindowSize, err = reader.ReadUInt32(binary.BigEndian)
	if err != nil {
		return nil, err
	}
	connect.ChannelCount, err = reader.ReadUInt32(binary.BigEndian)
	if err != nil {
		return nil, err
	}
	connect.IncomingBandwidth, err = reader.ReadUInt32(binary.BigEndian)
	if err != nil {
		return nil, err
	}
	connect.OutgoingBandwidth, err = reader.ReadUInt32(binary.BigEndian)
	if err != nil {
		return nil, err
	}
	connect.DisconnectThrottle, err = reader.ReadUInt32(binary.BigEndian)
	if err != nil {
		return nil, err
	}
	connect.PacketThrottleAcceleration, err = reader.ReadUInt32(binary.BigEndian)
	if err != nil {
		return nil, err
	}
	connect.PacketThrottleDeceleration, err = reader.ReadUInt32(binary.BigEndian)
	if err != nil {
		return nil, err
	}
	return &connect, nil
}
