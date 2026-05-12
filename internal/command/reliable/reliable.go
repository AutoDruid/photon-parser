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

// Reliable represents a complete reliable message with header and parameters.
// Parameters contain the actual game data as key-value pairs where each
// parameter has an ID, type, and value.
type Reliable[P types.ParameterView] struct {
	types.Reliable[P]
}

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
func Parse[P types.ParameterView](ctx *context.Context[P], length uint32) (*Reliable[P], error) {
	reliable := Reliable[P]{}
	header, err := reliable.parseHeader(ctx, length)
	if err != nil {
		return nil, err
	}

	if header.Type >= types.ExchangeKeys {
		return nil, nil
	}

	reliable.ReliableHeader = header

	if reliable.Signature != 0xF3 {
		return nil, errors.ErrEncryptedPacket
	}

	items := ctx.PoolParameter.Get(header.ParameterCount)
	reliable.Parameters = items.Items
	defer ctx.PoolParameter.Put(items)

	for i := 0; i < reliable.ParameterCount; i++ {
		err := ctx.Decoders.ParameterParser.Parse(ctx.Reader, &reliable.Parameters[i], ctx.Hooks)
		if err != nil {
			return nil, err
		}
	}

	emit(ctx.Hooks, &reliable)

	return &reliable, nil

}

func emit[P types.ParameterView](hooks *hooks.Hooks[P], out *Reliable[P]) {
	if hooks == nil {
		return
	}

	if hooks.OnEvents[out.Type] != nil {
		hooks.OnEvents[out.Type](out.Reliable)
	}
}

func (r *Reliable[P]) parseHeader(ctx *context.Context[P], length uint32) (types.ReliableHeader, error) {
	var err error
	var header types.ReliableHeader

	header.Signature, err = ctx.Reader.ReadUInt8()
	if err != nil {
		return types.ReliableHeader{}, err
	}

	b, err := ctx.Reader.ReadUInt8()
	if err != nil {
		return types.ReliableHeader{}, err
	}

	header.Type = types.Type(b)

	switch header.Type {
	case types.OperationResponse, types.OtherOperationResponse:

		header.EventCode, err = ctx.Reader.ReadUInt8()
		if err != nil {
			return types.ReliableHeader{}, err
		}

		//Return code
		_, err = ctx.Reader.ReadInt16(binary.LittleEndian)
		if err != nil {
			return types.ReliableHeader{}, err
		}

		//Read debug msg
		_, err = ctx.Reader.ReadByte()
		if err != nil {
			return types.ReliableHeader{}, err
		}
	case types.EventDataType, types.OperationRequest:
		header.EventCode, err = ctx.Reader.ReadUInt8()
		if err != nil {
			return types.ReliableHeader{}, err
		}
	default:
		_, err = ctx.Reader.ReadBytes(int(length) - READED_HEADER_SIZE)
		if err != nil {
			return types.ReliableHeader{}, err
		}
		return header, nil
	}

	header.ParameterCount, err = ctx.Decoders.ReliableHeaderParameterCount.Count(ctx.Reader)
	if err != nil {
		return types.ReliableHeader{}, err
	}

	return header, nil
}
