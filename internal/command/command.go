package command

import (
	"encoding/binary"
	"fmt"
	"michelprogram/photon-parser/internal/command/acknowledge"
	"michelprogram/photon-parser/internal/command/connect"
	"michelprogram/photon-parser/internal/command/disconnect"
	"michelprogram/photon-parser/internal/command/ping"
	"michelprogram/photon-parser/internal/command/reliable"
	"michelprogram/photon-parser/internal/context"
	"michelprogram/photon-parser/internal/errors"
	"michelprogram/photon-parser/internal/hooks"
	"michelprogram/photon-parser/internal/reader"
	"michelprogram/photon-parser/internal/types"
)

type Command struct {
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
func Parse(ctx *context.Context, out *types.Command) error {
	cmd := Command{}
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
		return errors.HeaderSize
	}

	parsed, err := cmd.parsePayload(header.Type, ctx, header.Length)
	if err != nil {
		rest, _ := ctx.Reader.ReadBytes(int(header.Length - types.COMMAND_HEADER_SIZE))
		// don't fatal — just store raw for encrypted packets
		out.Payload = types.UnknownPayload{Raw: rest, Kind: header.Type}
	} else {
		out.Payload = parsed
	}

	cmd.emit(ctx.Reader, ctx.Hooks)

	return nil
}

func (c Command) emit(r *reader.Reader, hooks *hooks.Hooks) {
	if hooks == nil {
		return
	}

	if hooks.SyncHooks.OnCommand != nil {
		hooks.SyncHooks.OnCommand(c.Command)
	}

	if hooks.AsyncHooks.OnCommand == nil {
		return
	}

	select {
	case hooks.AsyncHooks.OnCommand <- c.Command:
	default: // don't block parser
	}

}

func (c Command) parsePayload(t types.CommandType, ctx *context.Context, length uint32) (types.Payload, error) {
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
		return acknowledge.Parse(ctx)
	case types.ConnectCommand, types.VerifyConnectCommand:
		return connect.Parse(ctx)
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

func (s *Command) parseHeader(r *reader.Reader) (types.CommandHeader, error) {
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

func (c Command) String() string {
	return fmt.Sprintf("Type: %d, ChannelID: %d, Flags: %d, ReservedByte: %d, Length: %d, ReliableSequenceNumber: %d", c.Type, c.ChannelID, c.Flags, c.ReservedByte, c.Length, c.ReliableSequenceNumber)
}
