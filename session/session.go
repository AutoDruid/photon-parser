package session

import (
	"michelprogram/photon-parser/command"
	"michelprogram/photon-parser/parser"
)

// Parse parses a complete Photon session packet from a byte slice.
// This is a convenience wrapper around ParseFromReader.
//
// The session packet format consists of:
//   - Session header (12 bytes: peer ID, CRC flag, command count, timestamp, challenge)
//   - Variable number of commands (as specified in CommandCount)
//
// Returns a Session struct containing the parsed header and commands,
// or an error if parsing fails.
//
// Example usage:
//
//	session, err := session.Parse(packetBytes)
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Session has %d commands\n", len(session.Commands))
func Parse(data []byte) (*Session, error) {
	return ParseFromReader(parser.NewReader(data))
}

// ParseFromReader parses a Photon session packet from a parser.Reader.
// This function reads the session header, then iterates through and parses
// each command as specified by the CommandCount field.
//
// Returns a Session struct with all fields populated including the Commands slice,
// or an error if any part of parsing fails.
func ParseFromReader(r *parser.Reader) (*Session, error) {
	res := Session{}

	header, err := parser.ReadHeader[Header](r)
	if err != nil {
		return nil, err
	}

	res.Commands = make([]*command.Command, header.CommandCount)

	for i := uint8(0); i < header.CommandCount; i++ {
		cmd, err := command.ParseFromReader(r)
		if err != nil {
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
