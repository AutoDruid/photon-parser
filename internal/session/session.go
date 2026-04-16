// Package session provides parsing for Photon Protocol session layer packets.
// The session layer is the outermost protocol layer, containing session metadata
// and one or more commands.
package session

import (
	"fmt"
	"michelprogram/photon-parser/internal/command"
	"michelprogram/photon-parser/internal/reader"
	"michelprogram/photon-parser/internal/types"
)

type Session struct {
	types.Session
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

	s.Commands = make([]*types.Command, header.CommandCount)

	for i := uint8(0); i < header.CommandCount; i++ {
		cmd := &command.Command{}
		err := cmd.Parse(r)
		if err != nil {
			return err
		}
		s.Commands[i] = &cmd.Command
	}

	s.Header = header

	s.emit(r)

	return nil
}

func (s Session) emit(r *reader.Reader) {

	if r.SyncHooks.OnSession != nil {
		r.SyncHooks.OnSession(s.Session)
	}

	if r.AsyncHooks.OnSession == nil {
		return
	}

	select {
	case r.AsyncHooks.OnSession <- s.Session:
	default:
	}
}

func (s *Session) parseHeader(r *reader.Reader) (types.Header, error) {
	var err error
	var header types.Header

	header.PeerID, err = r.ReadUInt16()
	if err != nil {
		return types.Header{}, err
	}

	header.CRCEnabled, err = r.ReadUInt8()
	if err != nil {
		return types.Header{}, err
	}

	header.CommandCount, err = r.ReadUInt8()
	if err != nil {
		return types.Header{}, err
	}

	header.Timestamp, err = r.ReadUInt32()
	if err != nil {
		return types.Header{}, err
	}

	header.Challenge, err = r.ReadInt32()
	if err != nil {
		return types.Header{}, err
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
