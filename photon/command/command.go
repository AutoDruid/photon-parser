package command

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

func Parse(reader *bytes.Reader) (*Command, error){
	var header Header
	if err := binary.Read(reader, binary.BigEndian, &header); err != nil {
		return nil, err
	}

	if header.Length < HEADER_SIZE {
		return nil, fmt.Errorf("command length %d smaller than header size", header.Length)
	}

	dataLen := int(header.Length - HEADER_SIZE)
	payload := make([]byte, dataLen)

	if _, err := io.ReadFull(reader, payload); err != nil {
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

	return cmd,nil
}
