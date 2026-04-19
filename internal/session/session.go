// Package session provides parsing for Photon Protocol session layer packets.
// The session layer is the outermost protocol layer, containing session metadata
// and one or more commands.
package session

import (
	"fmt"
	"michelprogram/photon-parser/internal/command"
	"michelprogram/photon-parser/internal/hooks"
	"michelprogram/photon-parser/internal/reader"
	"michelprogram/photon-parser/internal/types"
)

type Session struct {
	types.Session
}

// Parse parses a Photon session packet from a parser.Reader.
// This function reads the session header, then iterates through and parses
// each command as specified by the CommandCount field.
//
// Returns a Session struct with all fields populated including the Commands slice,
// or an error if any part of parsing fails.
func Parse(reader *reader.Reader, hooks *hooks.Hooks) (*Session, error) {
	session := Session{}
	header, err := session.parseHeader(reader)
	if err != nil {
		return nil, err
	}

	session.Commands = make([]*types.Command, header.CommandCount)

	for i := uint8(0); i < header.CommandCount; i++ {
		cmd, err := command.Parse(reader, hooks)
		if err != nil {
			return nil, err
		}
		session.Commands[i] = &cmd.Command
	}

	session.Header = header

	session.emit(reader, hooks)

	return &session, nil
}

func (s Session) emit(reader *reader.Reader, hooks *hooks.Hooks) {
	if hooks == nil {
		return
	}

	if hooks.SyncHooks.OnSession != nil {
		hooks.SyncHooks.OnSession(s.Session)
	}

	if hooks.AsyncHooks.OnSession == nil {
		return
	}

	select {
	case hooks.AsyncHooks.OnSession <- s.Session:
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
