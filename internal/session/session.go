// Package session provides parsing for Photon Protocol session layer packets.
// The session layer is the outermost protocol layer, containing session metadata
// and one or more commands.
package session

import (
	"fmt"
	"michelprogram/photon-parser/internal/command"
	"michelprogram/photon-parser/internal/reader"
)

// Header represents the Photon session header containing peer and timing information.
// This header appears at the start of every Photon packet.
type Header struct {
	PeerID       uint16 `json:"peer_id"`       // Peer identifier for this connection
	CRCEnabled   uint8  `json:"crc_enabled"`   // CRC checksum flag (0 = disabled, 1 = enabled)
	CommandCount uint8  `json:"command_count"` // Number of commands following this header
	Timestamp    uint32 `json:"timestamp"`     // Timestamp in milliseconds
	Challenge    int32  `json:"challenge"`     // Challenge value for connection verification
}

// Session represents a complete Photon session packet with its header and commands.
// A session packet can contain multiple commands that will be processed sequentially.
type Session struct {
	Header

	Commands []*command.Command // Slice of commands contained in this session
}

var _ reader.Parseable = (*Session)(nil)

// Parse parses a Photon session packet from a parser.Reader.
// This function reads the session header, then iterates through and parses
// each command as specified by the CommandCount field.
//
// Returns a Session struct with all fields populated including the Commands slice,
// or an error if any part of parsing fails.
func (s *Session) Parse(r *reader.Reader) error {
	header, err := s.parseHeader(r)
	if err != nil {
		return err
	}

	s.Commands = make([]*command.Command, header.CommandCount)

	for i := uint8(0); i < header.CommandCount; i++ {
		cmd := &command.Command{}
		err := cmd.Parse(r)
		if err != nil {
			return err
		}
		s.Commands[i] = cmd
	}

	s.Header = header

	return nil
}

func (s *Session) parseHeader(r *reader.Reader) (Header, error) {
	var err error
	var header Header

	header.PeerID, err = r.ReadUInt16()
	if err != nil {
		return Header{}, err
	}

	header.CRCEnabled, err = r.ReadUInt8()
	if err != nil {
		return Header{}, err
	}

	header.CommandCount, err = r.ReadUInt8()
	if err != nil {
		return Header{}, err
	}

	header.Timestamp, err = r.ReadUInt32()
	if err != nil {
		return Header{}, err
	}

	header.Challenge, err = r.ReadInt32()
	if err != nil {
		return Header{}, err
	}

	return header, nil
}

func (s Session) String() string {
	res := fmt.Sprintf("Session: PeerID: %d, CRCEnabled: %d, CommandCount: %d, Timestamp: %d, Challenge: %d", s.PeerID, s.CRCEnabled, s.CommandCount, s.Timestamp, s.Challenge)
	for i, cmd := range s.Commands {
		res += fmt.Sprintf("\n  Command %d: Type: %d, Payload: %v", i, cmd.Type, cmd.Payload)
	}
	return res
}
