package session

import "michelprogram/photon-parser/photon/command"

type Header struct{
	PeerID       uint16
	CRCEnabled   uint8
	CommandCount uint8
	Timestamp    uint32
	Challenge    int32
}

type Session struct {
	Header

	Commands []*command.Command
}
