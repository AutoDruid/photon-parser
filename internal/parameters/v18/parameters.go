package v18

import (
	"fmt"
	"log"
	"michelprogram/photon-parser/internal/hooks"
	"michelprogram/photon-parser/internal/reader"
	"michelprogram/photon-parser/internal/types"
)

type Parameter struct {
	types.Parameter
}

var _ reader.ParameterParser = (*Parameter)(nil)

// Parse reads a complete parameter from the reader.
// Format: Header (1 byte ID + 1 byte Type), followed by the typed value.
//
// The function first reads the parameter header to determine the parameter ID
// and type code, then decodes the value according to that type using the
// Protocol16 decoder.
//
// Returns a Parameters struct containing the ID, Type, and decoded Value,
// or an error if parsing fails.
//
// Example usage:
//
//	param, err := parameters.Parse(reader)
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Parameter %d has value: %v\n", param.ID, param.Value)
func (p *Parameter) Parse(reader *reader.Reader, out *types.Parameter, hooks *hooks.Hooks) error {

	header, err := p.parseHeader(reader)
	if err != nil {
		return err
	}

	log.Println("header", header)

	value, err := p.decode(reader, ParameterType(header.Type))

	log.Println("value", value)
	if reader.Cursor+1 < reader.Max {
		log.Printf("next % x\n", reader.Buffer[reader.Cursor+1:])

	}

	if err != nil {
		log.Println("err on parameter type", header.Type, err)
		return err

	}

	out.ParameterHeader = header
	out.Value = value

	p.emit(reader, hooks)

	return nil
}

func (p Parameter) emit(reader *reader.Reader, hooks *hooks.Hooks) {
	if hooks == nil {
		return
	}

	if hooks.SyncHooks.OnParameter != nil {
		hooks.SyncHooks.OnParameter(p.Parameter)
	}

	if hooks.AsyncHooks.OnParameter == nil {
		return
	}

	select {
	case hooks.AsyncHooks.OnParameter <- p.Parameter:
	default:
	}
}

func (p *Parameter) parseHeader(r *reader.Reader) (types.ParameterHeader, error) {
	var err error
	var header types.ParameterHeader

	header.ID, err = r.ReadUInt8()
	if err != nil {
		return types.ParameterHeader{}, err
	}

	b, err := r.ReadUInt8()
	if err != nil {
		return types.ParameterHeader{}, err
	}

	header.Type = types.ParameterType(b)

	return header, nil
}

// Decode reads a value of the specified Photon Protocol16 type from the reader.
// It dispatches to the appropriate type-specific reader based on ttype.
// Returns the decoded value as any, or an error if the type is unsupported
// or if reading fails.
//
// Supported types include all primitives (int8, int16, int32, int64, float32, float64,
// string, boolean), arrays, dictionaries, and hashtables.
//
// For NilType and UnknownType, returns nil without error.
// For unsupported type codes, returns an error.
func (p Parameter) decode(reader *reader.Reader, ttype ParameterType) (any, error) {
	switch ttype {
	case Int16Type:
		return reader.ReadInt16LittleEndian()
	case Int16Positive:
		value, err := reader.ReadUInt16LittleEndian()
		if err != nil {
			return nil, err
		}
		return int32(value), nil
	case Long16Positive:
		value, err := reader.ReadUInt16LittleEndian()
		if err != nil {
			return nil, err
		}
		return int64(value), nil
	case StringType:
		return p.readString(reader)
	case CompressedInt32Type:
		return reader.ReadVarintInt32()
	case CompressedInt64Type:
		return reader.ReadVarintInt64()
	case Float32ArrayType:
		return p.readFloatArray(reader)
	case Float32Type:
		return reader.ReadFloat32()
	case Int8Type:
		return reader.ReadInt8()
	case ByteZeroType:
		return byte(0), nil
	case ShortArrayType:
		return p.readInt16Array(reader)
	case ByteArrayType:
		return p.readInt8Array(reader)
	case NilType, UnknownType:
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported type: %d", ttype)
	}
}

// String returns a human-readable representation of the parameter.
// Format: "ID: <id>\nType: <type>\nValue: <value>\n"
func (p Parameter) String() string {
	return fmt.Sprintf("ID: %d\nType: %d\nValue: %v\n", p.ID, p.Type, p.Value)
}
