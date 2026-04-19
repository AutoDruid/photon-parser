package v16

import (
	"fmt"
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

	value, err := p.decode(reader, header.Type)

	if err != nil {
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
func (p Parameter) decode(reader *reader.Reader, ttype types.ParameterType) (any, error) {
	switch ttype {
	case types.Int8Type:
		return reader.ReadInt8()
	case types.Int16Type:
		return reader.ReadInt16()
	case types.Int32Type:
		return reader.ReadInt32()
	case types.Int64Type:
		return reader.ReadInt64()
	case types.Float32Type:
		return reader.ReadFloat32()
	case types.Float64Type:
		return reader.ReadFloat64()
	case types.StringType:
		return p.readString(reader)
	case types.BooleanType:
		return reader.ReadBoolean()
	case types.Int8ArrayType:
		return p.readInt8Array(reader)
	case types.Int32ArrayType:
		return p.readInt32Array(reader)
	case types.ArrayType:
		return p.readArray(reader)
	case types.StringArrayType:
		return p.readStringArray(reader)
	case types.DictionaryType:
		return p.readDictionary(reader)
	case types.HashTableType:
		return p.readHashTable(reader)
	case types.NilType, types.UnknownType:
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported type: 0x%02x", ttype)
	}
}

// String returns a human-readable representation of the parameter.
// Format: "ID: <id>\nType: <type>\nValue: <value>\n"
func (p Parameter) String() string {
	return fmt.Sprintf("ID: %d\nType: %d\nValue: %v\n", p.ID, p.Type, p.Value)
}
