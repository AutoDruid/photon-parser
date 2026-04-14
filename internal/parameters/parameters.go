package parameters

import (
	"fmt"
	"michelprogram/photon-parser/internal/reader"
)

// Header represents the parameter header containing the parameter ID and type.
// This appears at the beginning of each serialized parameter.
type Header struct {
	ID   uint8 // Parameter identifier (application-specific)
	Type Type  // Protocol16 type code indicating how to decode the value
}

// Parameters represents a complete Photon Protocol parameter with its header and decoded value.
// The Value field contains the decoded data according to the Type specified in the Header.
type Parameters struct {
	Header

	Value interface{} // Decoded value, type depends on Header.Type
}

var _ reader.Parseable = (*Parameters)(nil)

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
func (p *Parameters) Parse(r *reader.Reader) error {

	header, err := p.parseHeader(r)
	if err != nil {
		return err
	}

	value, err := p.decode(r, header.Type)

	if err != nil {
		return err
	}

	p.Header = header
	p.Value = value

	return nil
}

func (p *Parameters) parseHeader(r *reader.Reader) (Header, error) {
	var err error
	var header Header

	header.ID, err = r.ReadUInt8()
	if err != nil {
		return Header{}, err
	}

	b, err := r.ReadUInt8()
	if err != nil {
		return Header{}, err
	}

	header.Type = Type(b)

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
func (p Parameters) decode(reader *reader.Reader, ttype Type) (any, error) {
	switch ttype {
	case Int8Type:
		return reader.ReadInt8()
	case Int16Type:
		return reader.ReadInt16()
	case Int32Type:
		return reader.ReadInt32()
	case Int64Type:
		return reader.ReadInt64()
	case Float32Type:
		return reader.ReadFloat32()
	case Float64Type:
		return reader.ReadFloat64()
	case StringType:
		return p.readString(reader)
	case BooleanType:
		return reader.ReadBoolean()
	case Int8ArrayType:
		return p.readInt8Array(reader)
	case Int32ArrayType:
		return p.readInt32Array(reader)
	case ArrayType:
		return p.readArray(reader)
	case StringArrayType:
		return p.readStringArray(reader)
	case DictionaryType:
		return p.readDictionary(reader)
	case HashTableType:
		return p.readHashTable(reader)
	case NilType, UnknownType:
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported type: 0x%02x", ttype)
	}
}

// String returns a human-readable representation of the parameter.
// Format: "ID: <id>\nType: <type>\nValue: <value>\n"
func (p Parameters) String() string {
	return fmt.Sprintf("ID: %d\nType: %d\nValue: %v\n", p.ID, p.Type, p.Value)
}
