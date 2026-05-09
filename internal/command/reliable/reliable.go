package reliable

import (
	"encoding/binary"
	"michelprogram/photon-parser/internal/context"
	"michelprogram/photon-parser/internal/errors"
	"michelprogram/photon-parser/internal/types"
)

// HEADER_SIZE is the size in bytes of a reliable message header (5 bytes).
const HEADER_SIZE = 5

// READED_HEADER_SIZE is the number of bytes consumed before parameter parsing in reliable commands.
const READED_HEADER_SIZE = 14

// Type represents a Photon reliable message type.
type Type uint8

// Photon Protocol reliable message types.
// These define the different kinds of reliable messages that can be exchanged.
const (
	OperationRequest       Type = 0x02 // Client requests an operation
	OperationResponse      Type = 0x07 // Server responds to an operation
	OtherOperationResponse Type = 0x03 // Alternative response format
	EventDataType          Type = 0x04 // Server sends an event to client
	ExchangeKeys           Type = 0x06 // Key exchange for encryption
)

// Header represents the reliable message header.
// This appears at the start of the payload in SendReliable commands.
type Header struct {
	Signature      uint8 `json:"signature"`       // Message signature (typically 0xF3)
	Type           Type  `json:"type"`            // Message type (operation, event, etc.)
	EventCode      uint8 `json:"event_code"`      // Operation/event code (application-specific)
	ParameterCount int   `json:"parameter_count"` // Number of parameters following this header
}

// Reliable represents a complete reliable message with header and parameters.
// Parameters contain the actual game data as key-value pairs where each
// parameter has an ID, type, and value.
type Reliable[P types.ParameterView] struct {
	Header
	Parameters []P // Slice of decoded parameters
}

// Parse parses a Photon reliable message from the parsing context.
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

	if header.Type >= ExchangeKeys {
		return nil, nil
	}

	reliable.Header = header

	if reliable.Signature != 0xF3 {
		return nil, errors.ErrEncryptedPacket
	}

	reliable.Parameters = make([]P, header.ParameterCount)

	for i := 0; i < reliable.ParameterCount; i++ {
		err := ctx.Decoders.ParameterParser.Parse(ctx.Reader, &reliable.Parameters[i], ctx.Hooks)
		if err != nil {
			return nil, err
		}
	}

	return &reliable, nil

}

func (r *Reliable[P]) parseHeader(ctx *context.Context[P], length uint32) (Header, error) {
	var err error
	var header Header

	header.Signature, err = ctx.Reader.ReadUInt8()
	if err != nil {
		return Header{}, err
	}

	b, err := ctx.Reader.ReadUInt8()
	if err != nil {
		return Header{}, err
	}

	header.Type = Type(b)

	switch header.Type {
	case OperationResponse, OtherOperationResponse:

		header.EventCode, err = ctx.Reader.ReadUInt8()
		if err != nil {
			return Header{}, err
		}

		//Return code
		_, err = ctx.Reader.ReadInt16(binary.LittleEndian)
		if err != nil {
			return Header{}, err
		}

		//Read debug msg
		_, err = ctx.Reader.ReadByte()
		if err != nil {
			return Header{}, err
		}
	case EventDataType, OperationRequest:
		header.EventCode, err = ctx.Reader.ReadUInt8()
		if err != nil {
			return Header{}, err
		}
	default:
		_, err = ctx.Reader.ReadBytes(int(length) - READED_HEADER_SIZE)
		if err != nil {
			return Header{}, err
		}
		return header, nil
	}

	header.ParameterCount, err = ctx.Decoders.ReliableHeaderParameterCount.Count(ctx.Reader)
	if err != nil {
		return Header{}, err
	}

	return header, nil
}
