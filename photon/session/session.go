package session

import (
	"bytes"
	"encoding/binary"
	"michelprogram/photon-parser/photon/command"
)


func Parse(packet []byte) (*Session, error) {
	res := Session{}
	reader := bytes.NewReader(packet)

	var header Header
	var i uint8 = 0

	err := binary.Read(reader, binary.BigEndian, &header)
	if err != nil {
		return nil, err
	}

	res.Commands = make([]*command.Command, header.CommandCount)

	for i = 0; i < header.CommandCount; i++ {
		cmd, err := command.Parse(reader)
		if err != nil{
			return nil, err
		}
		res.Commands[i] = cmd
	}
	
	res.PeerID = header.PeerID
	res.CommandCount = header.CommandCount
	res.Timestamp = header.Timestamp
	res.Challenge = header.Challenge
	res.CRCEnabled = header.CRCEnabled

	return &res, nil
}
