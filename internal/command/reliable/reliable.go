package reliable

import (
	"encoding/binary"

	"github.com/AutoDruid/photon-parser/internal/context"
	"github.com/AutoDruid/photon-parser/internal/errors"
	"github.com/AutoDruid/photon-parser/internal/hooks"
	"github.com/AutoDruid/photon-parser/internal/types"
)

// HEADER_SIZE is the size in bytes of a reliable message header (5 bytes).
const HEADER_SIZE = 5

const READED_HEADER_SIZE = 14

// ParseFromReader parses a Photon reliable message from a parser.Reader.
// It reads the 5-byte header, then iterates through and parses each parameter
// as specified by the ParameterCount field.
//
// The message format consists of:
//   - Header (5 bytes: signature, type, event code, parameter count)
//   - Parameters (ParameterCount entries, each with ID, type, and value)
//
// Returns a Reliable struct with all fields populated including the Parameters slice,
// or an error if any part of parsing fails.
func Parse[P types.ParameterView](ctx *context.Context[P], out *types.Reliable[P], length uint32) error {
	err := parseHeader(out, ctx, length)
	if err != nil {
		return err
	}

	if out.Type >= types.ExchangeKeys {
		return nil
	}

	if out.Signature != 0xF3 {
		return errors.ErrEncryptedPacket
	}

	items := ctx.PoolParameter.Get(out.ParameterCount)
	out.Parameters = items.Items
	defer ctx.PoolParameter.Put(items)

	for i := 0; i < out.ParameterCount; i++ {
		err := ctx.Decoders.ParameterParser.Parse(ctx.Reader, &out.Parameters[i], ctx.Hooks)
		if err != nil {
			return err
		}
	}

	emit(ctx.Hooks, out)

	return nil

}

func emit[P types.ParameterView](hooks *hooks.Hooks[P], out *types.Reliable[P]) {
	if hooks == nil {
		return
	}

	if hooks.OnEvents[out.Type] != nil {
		hooks.OnEvents[out.Type](*out)
	}
}

func parseHeader[P types.ParameterView](out *types.Reliable[P], ctx *context.Context[P], length uint32) error {
	var err error

	out.Signature, err = ctx.Reader.ReadUInt8()
	if err != nil {
		return err
	}

	b, err := ctx.Reader.ReadUInt8()
	if err != nil {
		return err
	}

	out.Type = types.Type(b)

	switch out.Type {
	case types.OperationResponse, types.OtherOperationResponse:

		out.EventCode, err = ctx.Reader.ReadUInt8()
		if err != nil {
			return err
		}

		//Return code
		_, err = ctx.Reader.ReadInt16(binary.LittleEndian)
		if err != nil {
			return err
		}

		//Read debug msg
		_, err = ctx.Reader.ReadByte()
		if err != nil {
			return err
		}
	case types.EventDataType, types.OperationRequest:
		out.EventCode, err = ctx.Reader.ReadUInt8()
		if err != nil {
			return err
		}
	default:
		_, err = ctx.Reader.ReadBytes(int(length) - READED_HEADER_SIZE)
		if err != nil {
			return err
		}
		return nil
	}

	out.ParameterCount, err = ctx.Decoders.ReliableHeaderParameterCount.Count(ctx.Reader)
	if err != nil {
		return err
	}

	return nil
}
