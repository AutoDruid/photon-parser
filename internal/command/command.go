package command

import (
	"fmt"

	"github.com/AutoDruid/photon-parser/internal/command/acknowledge"
	"github.com/AutoDruid/photon-parser/internal/command/connect"
	"github.com/AutoDruid/photon-parser/internal/command/reliable"
	"github.com/AutoDruid/photon-parser/internal/context"
	"github.com/AutoDruid/photon-parser/internal/errors"
	"github.com/AutoDruid/photon-parser/internal/hooks"
	"github.com/AutoDruid/photon-parser/internal/reader"
	"github.com/AutoDruid/photon-parser/internal/types"
)

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
func Parse[P types.ParameterView](ctx *context.Context[P], out *types.Command[P]) error {
	err := parseHeader(out, ctx.Reader)

	if out.Type > types.SendReliableFragmentCommand {
		remaining := ctx.Reader.Max - ctx.Reader.Cursor - 1
		rest, err := ctx.Reader.ReadBytes(remaining)
		if err != nil {
			return err
		}
		out.UnknownPayload = types.UnknownPayload{Raw: rest, Kind: out.Type}
		return nil
	}

	if err != nil {
		return err
	}

	if out.Length < types.COMMAND_HEADER_SIZE {
		return errors.ErrHeaderSize
	}

	err = parsePayload(out, ctx)
	if err != nil {
		rest, _ := ctx.Reader.ReadBytes(int(out.Length - types.COMMAND_HEADER_SIZE))
		// don't fatal — just store raw for encrypted packets
		out.UnknownPayload = types.UnknownPayload{Raw: rest, Kind: out.Type}
	}

	emit(ctx.Hooks, out)

	return nil
}

func emit[P types.ParameterView](hooks *hooks.Hooks[P], out *types.Command[P]) {
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

func parsePayload[P types.ParameterView](out *types.Command[P], ctx *context.Context[P]) error {
	switch out.Type {
	case types.SendUnreliableCommand:

		_, err := ctx.Reader.ReadBytes(4)
		if err != nil {
			return err
		}
		err = reliable.Parse(ctx, &out.UnreliablePayload, out.Length)
		if err != nil {
			return err
		}
	case types.SendReliableCommand:
		err := reliable.Parse(ctx, &out.ReliablePayload, out.Length)
		if err != nil {
			return err
		}
	case types.AcknowledgeCommand:
		err := acknowledge.Parse(ctx.Reader, &out.AcknowledgePayload)
		if err != nil {
			return err
		}
	case types.ConnectCommand, types.VerifyConnectCommand:
		err := connect.Parse(ctx.Reader, &out.ConnectPayload)
		if err != nil {
			return err
		}
	case types.SendReliableFragmentCommand:
		err := reliable.ParseFragment(ctx, &out.ReliableFragmentPayload, &out.ReliablePayload, out.Length)
		if err != nil {
			return err
		}
	case types.PingCommand:
		out.PingPayload = struct{}{}
	case types.DisconnectCommand:
		out.DisconnectPayload = struct{}{}
	default:
		return fmt.Errorf("unsupported command type %d", out.Type)
	}

	return nil
}

func parseHeader[P types.ParameterView](out *types.Command[P], r *reader.Reader) error {
	var err error

	b, err := r.ReadUInt8()
	if err != nil {
		return err
	}

	out.Type = types.CommandType(b)

	if out.Type > types.SendReliableFragmentCommand {
		return nil
	}

	out.ChannelID, err = r.ReadUInt8()
	if err != nil {
		return err
	}

	out.Flags, err = r.ReadUInt8()
	if err != nil {
		return err
	}

	out.ReservedByte, err = r.ReadUInt8()
	if err != nil {
		return err
	}

	out.Length, err = r.ReadUInt32BE()
	if err != nil {
		return err
	}

	out.ReliableSequenceNumber, err = r.ReadUInt32BE()
	if err != nil {
		return err
	}

	return nil
}
