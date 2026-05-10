package command

import (
	"encoding/binary"
	"fmt"

	"github.com/AutoDruid/photon-parser/internal/command/acknowledge"
	"github.com/AutoDruid/photon-parser/internal/command/connect"
	"github.com/AutoDruid/photon-parser/internal/command/disconnect"
	"github.com/AutoDruid/photon-parser/internal/command/ping"
	"github.com/AutoDruid/photon-parser/internal/command/reliable"
	"github.com/AutoDruid/photon-parser/internal/context"
	"github.com/AutoDruid/photon-parser/internal/errors"
	"github.com/AutoDruid/photon-parser/internal/hooks"
	"github.com/AutoDruid/photon-parser/internal/reader"
	"github.com/AutoDruid/photon-parser/internal/types"
)

type Command[P types.ParameterView] struct {
	types.Command
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
func Parse[P types.ParameterView](ctx *context.Context[P], out *types.Command) error {
	cmd := Command[P]{}
	header, err := cmd.parseHeader(ctx.Reader)

	out.CommandHeader = header
	if header.Type > types.SendReliableFragmentCommand {
		remaining := ctx.Reader.Max - ctx.Reader.Cursor - 1
		rest, err := ctx.Reader.ReadBytes(remaining)
		if err != nil {
			return err
		}
		out.Payload = types.UnknownPayload{Raw: rest, Kind: header.Type}
		return nil
	}

	if err != nil {
		return err
	}

	if header.Length < types.COMMAND_HEADER_SIZE {
		return errors.ErrHeaderSize
	}

	parsed, err := cmd.parsePayload(header.Type, ctx, header.Length)
	if err != nil {
		rest, _ := ctx.Reader.ReadBytes(int(header.Length - types.COMMAND_HEADER_SIZE))
		// don't fatal — just store raw for encrypted packets
		out.Payload = types.UnknownPayload{Raw: rest, Kind: header.Type}
	} else {
		out.Payload = parsed
	}

	cmd.emit(ctx.Hooks, out)

	return nil
}

func (c Command[P]) emit(hooks *hooks.Hooks[P], out *types.Command) {
	if hooks == nil {
		return
	}

	if hooks.SyncHooks.OnCommand != nil {
		hooks.SyncHooks.OnCommand(*out)
	}

	if hooks.AsyncHooks.OnCommand == nil {
		return
	}

	select {
	case hooks.AsyncHooks.OnCommand <- *out:
	default: // don't block parser
	}

}

func (c Command[P]) parsePayload(t types.CommandType, ctx *context.Context[P], length uint32) (types.Payload, error) {
	switch t {
	case types.SendUnreliableCommand:

		_, err := ctx.Reader.ReadBytes(4)
		if err != nil {
			return nil, err
		}
		sd, err := reliable.Parse(ctx, length)
		if err != nil {
			return nil, err
		}
		return sd, nil
	case types.SendReliableCommand:
		return reliable.Parse(ctx, length)
	case types.AcknowledgeCommand:
		return acknowledge.Parse(ctx.Reader)
	case types.ConnectCommand, types.VerifyConnectCommand:
		return connect.Parse(ctx.Reader)
	case types.SendReliableFragmentCommand:
		return reliable.ParseFragment(ctx, length)
	case types.PingCommand:
		return ping.Parse(), nil
	case types.DisconnectCommand:
		return disconnect.Parse(), nil
	default:
		return nil, fmt.Errorf("unsupported command type %d", t)
	}
}

func (s *Command[P]) parseHeader(r *reader.Reader) (types.CommandHeader, error) {
	var err error
	var header types.CommandHeader

	b, err := r.ReadUInt8()
	if err != nil {
		return types.CommandHeader{}, err
	}

	header.Type = types.CommandType(b)

	if header.Type > types.SendReliableFragmentCommand {
		return header, nil
	}

	header.ChannelID, err = r.ReadUInt8()
	if err != nil {
		return types.CommandHeader{}, err
	}

	header.Flags, err = r.ReadUInt8()
	if err != nil {
		return types.CommandHeader{}, err
	}

	header.ReservedByte, err = r.ReadUInt8()
	if err != nil {
		return types.CommandHeader{}, err
	}

	header.Length, err = r.ReadUInt32(binary.BigEndian)
	if err != nil {
		return types.CommandHeader{}, err
	}

	header.ReliableSequenceNumber, err = r.ReadUInt32(binary.BigEndian)
	if err != nil {
		return types.CommandHeader{}, err
	}

	return header, nil
}
