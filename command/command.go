package command

import (
	"fmt"
	"michelprogram/photon-parser/parser"
)

func Parse(data []byte) (*Command, error) {
	return ParseFromReader(parser.NewReader(data))
}

func ParseFromReader(r *parser.Reader) (*Command, error) {
	header, err := parser.ReadHeader[Header](r)
	if err != nil {
		return nil, err
	}

	if header.Length < HEADER_SIZE {
		return nil, fmt.Errorf("command length %d smaller than header size", header.Length)
	}

	payload, err := r.ReadBytes(int(header.Length - HEADER_SIZE))

	if err != nil {
		return nil, err
	}

	cmd := &Command{}

	cmd.Type = header.Type
	cmd.ChannelID = header.ChannelID
	cmd.Flags = header.Flags
	cmd.ReservedByte = header.ReservedByte
	cmd.Length = header.Length
	cmd.ReliableSequenceNumber = header.ReliableSequenceNumber
	cmd.Data = payload

	return cmd, nil
}
