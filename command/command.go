package command

import (
	"fmt"
	"michelprogram/photon-parser/parser"
)

// Parse parses a Photon command from a byte slice.
// This is a convenience wrapper around ParseFromReader.
//
// The command format consists of:
//   - Command header (12 bytes: type, channel, flags, reserved, length, sequence number)
//   - Command payload (length - 12 bytes)
//
// Returns a Command struct containing the parsed header and payload data,
// or an error if parsing fails.
//
// Example usage:
//
//	cmd, err := command.Parse(commandBytes)
//	if err != nil {
//	    return err
//	}
//	if cmd.Type == command.SendReliable {
//	    // Process reliable message
//	}
func Parse(data []byte) (*Command, error) {
	return ParseFromReader(parser.NewReader(data))
}

// ParseFromReader parses a Photon command from a parser.Reader.
// It first reads the 12-byte command header, validates the length field,
// then reads the remaining payload data.
//
// Returns an error if:
//   - The header cannot be read
//   - The length field is smaller than the header size (invalid)
//   - The payload data cannot be fully read
//
// The returned Command struct contains all header fields and the raw payload
// in the Data field. For SendReliable commands, the Data can be further parsed
// using the reliable package.
func ParseFromReader(r *parser.Reader) (*Command, error) {
	header, err := parser.ReadHeader[Header](r)
	if err != nil {
		return nil, err
	}

	if header.Length < HEADER_SIZE {
		return nil, fmt.Errorf("command length %d smaller than header size %d", header.Length, HEADER_SIZE)
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
