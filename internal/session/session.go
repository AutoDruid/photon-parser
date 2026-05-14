// Package session provides parsing for Photon Protocol session layer packets.
// The session layer is the outermost protocol layer, containing session metadata
// and one or more commands.
package session

import (
	"errors"

	"github.com/AutoDruid/photon-parser/internal/command"
	"github.com/AutoDruid/photon-parser/internal/context"
	photonErrors "github.com/AutoDruid/photon-parser/internal/errors"
	"github.com/AutoDruid/photon-parser/internal/hooks"
	"github.com/AutoDruid/photon-parser/internal/reader"
	"github.com/AutoDruid/photon-parser/internal/types"
)

type Session[P types.ParameterView] struct {
	types.Session[P]
}

// Parse parses a Photon session packet from a parser.Reader.
// This function reads the session header, then iterates through and parses
// each command as specified by the CommandCount field.
//
// Returns a Session struct with all fields populated including the Commands slice,
// or an error if any part of parsing fails.
func Parse[P types.ParameterView](ctx *context.Context[P], out *types.Session[P]) error {
	err := parseHeader(out, ctx.Reader)
	if err != nil {
		return err
	}

	items := ctx.PoolCommand.Get(int(out.CommandCount))
	out.Commands = items.Items

	for i := uint8(0); i < out.CommandCount; i++ {
		err := command.Parse(ctx, &out.Commands[i])

		if errors.Is(err, photonErrors.ErrHeaderSize) {
			break
		}
		if err != nil {
			return err
		}

		if out.Commands[i].Type > types.SendReliableFragmentCommand {
			break
		}
	}

	emit(ctx.Hooks, out)

	ctx.PoolCommand.Put(items)

	return nil
}

func parseHeader[P types.ParameterView](out *types.Session[P], r *reader.Reader) error {
	var err error

	out.PeerID, err = r.ReadUInt16BE()
	if err != nil {
		return err
	}

	out.CRCEnabled, err = r.ReadUInt8()
	if err != nil {
		return err
	}

	out.CommandCount, err = r.ReadUInt8()
	if err != nil {
		return err
	}

	out.Timestamp, err = r.ReadUInt32BE()
	if err != nil {
		return err
	}

	out.Challenge, err = r.ReadInt32BE()
	if err != nil {
		return err
	}

	return nil
}

func emit[P types.ParameterView](hooks *hooks.Hooks[P], out *types.Session[P]) {
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
