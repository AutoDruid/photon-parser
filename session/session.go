package session

import (
	"michelprogram/photon-parser/command"
	"michelprogram/photon-parser/parser"
)

func Parse(data []byte) (*Session, error) {
    return ParseFromReader(parser.NewReader(data))
}

func ParseFromReader(r *parser.Reader) (*Session, error) {
	res := Session{}

	var i uint8 = 0

	header, err := parser.ReadHeader[Header](r)
	if err != nil {
		return nil, err
	}

	res.Commands = make([]*command.Command, header.CommandCount)

	for i = 0; i < header.CommandCount; i++ {
		cmd, err := command.ParseFromReader(r)
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
