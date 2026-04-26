// Package session provides parsing for Photon Protocol session layer packets.
// The session layer is the outermost protocol layer, containing session metadata
// and one or more commands.
package session

import (
	"encoding/binary"
	"errors"
	"michelprogram/photon-parser/internal/command"
	"michelprogram/photon-parser/internal/context"
	photonErrors "michelprogram/photon-parser/internal/errors"
	"michelprogram/photon-parser/internal/hooks"
	"michelprogram/photon-parser/internal/reader"
	"michelprogram/photon-parser/internal/types"
)

type Session[P types.ParameterView] struct {
	types.Session
}

// Parse parses a Photon session packet from a parser.Reader.
// This function reads the session header, then iterates through and parses
// each command as specified by the CommandCount field.
//
// Returns a Session struct with all fields populated including the Commands slice,
// or an error if any part of parsing fails.
func Parse[P types.ParameterView](ctx *context.Context[P], out *types.Session) error {
	session := Session[P]{}
	header, err := session.parseHeader(ctx.Reader)
	if err != nil {
		return err
	}

	out.Commands = make([]types.Command, header.CommandCount)

	for i := uint8(0); i < header.CommandCount; i++ {
		err := command.Parse(ctx, &out.Commands[i])

		if errors.Is(err, photonErrors.HeaderSize) {
			break
		}
		if err != nil {
			return err
		}

		if out.Commands[i].Type > types.SendReliableFragmentCommand {
			break
		}
	}

	out.Header = header

	session.emit(ctx.Hooks, out)

	return nil
}

func (s *Session[P]) parseHeader(r *reader.Reader) (types.Header, error) {
	var err error
	var header types.Header

	header.PeerID, err = r.ReadUInt16(binary.BigEndian)
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

	header.Timestamp, err = r.ReadUInt32(binary.BigEndian)
	if err != nil {
		return types.Header{}, err
	}

	header.Challenge, err = r.ReadInt32(binary.BigEndian)
	if err != nil {
		return types.Header{}, err
	}

	return header, nil
}

func (s Session[P]) emit(hooks *hooks.Hooks[P], out *types.Session) {
	if hooks == nil {
		return
	}

	if hooks.SyncHooks.OnSession != nil {
		hooks.SyncHooks.OnSession(*out)
	}

	if hooks.AsyncHooks.OnSession == nil {
		return
	}

	select {
	case hooks.AsyncHooks.OnSession <- *out:
	default:
	}
}
