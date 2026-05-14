package reliable

import (
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
func ParseInto[P types.ParameterView](ctx *context.Context[P], length uint32, dest *types.Reliable[P]) error {
	err := readReliableHeaderInto(ctx, length, dest)
	if err != nil {
		return err
	}

	if dest.Type >= types.ExchangeKeys {
		return nil
	}

	if dest.Signature != 0xF3 {
		return errors.ErrEncryptedPacket
	}

	items := ctx.PoolParameter.Get(dest.ParameterCount)
	dest.Parameters = items.Items

	for i := 0; i < dest.ParameterCount; i++ {
		err := ctx.Decoders.ParameterParser.ParseInto(ctx.Reader, ctx.Hooks, &dest.Parameters[i])
		if err != nil {
			return err
		}
	}

	emit(ctx.Hooks, dest)

	ctx.PoolParameter.Put(items)

	return nil

}

func readReliableHeaderInto[P types.ParameterView](ctx *context.Context[P], length uint32, dest *types.Reliable[P]) error {
	var err error

	dest.Signature, err = ctx.Reader.ReadUInt8()
	if err != nil {
		return err
	}

	b, err := ctx.Reader.ReadUInt8()
	if err != nil {
		return err
	}

	dest.Type = types.MessageType(b)

	switch dest.Type {
	case types.OperationResponse, types.OtherOperationResponse:

		dest.EventCode, err = ctx.Reader.ReadUInt8()
		if err != nil {
			return err
		}

		//Return code
		_, err = ctx.Reader.ReadInt16LE()
		if err != nil {
			return err
		}

		//Read debug msg
		_, err = ctx.Reader.ReadByte()
		if err != nil {
			return err
		}
	case types.EventDataType, types.OperationRequest:
		dest.EventCode, err = ctx.Reader.ReadUInt8()
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

	dest.ParameterCount, err = ctx.Decoders.ReliableHeaderParameterCount.Count(ctx.Reader)
	if err != nil {
		return err
	}

	return nil
}

func emit[P types.ParameterView](hooks *hooks.Hooks[P], dest *types.Reliable[P]) {
	if hooks == nil {
		return
	}

	if hooks.OnEvents[dest.Type] != nil {
		hooks.OnEvents[dest.Type](*dest)
	}
}
