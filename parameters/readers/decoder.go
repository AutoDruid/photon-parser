package readers

import (
	"fmt"
	"michelprogram/photon-parser/parser"
)

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
func Decode(reader *parser.Reader, ttype Type) (any, error) {
	switch ttype {
	case Int8Type:
		return ReadInt8(reader)
	case Int16Type:
		return ReadInt16(reader)
	case Int32Type:
		return ReadInt32(reader)
	case Int64Type:
		return ReadInt64(reader)
	case Float32Type:
		return ReadFloat32(reader)
	case Float64Type:
		return ReadFloat64(reader)
	case StringType:
		return ReadString(reader)
	case BooleanType:
		return ReadBoolean(reader)
	case Int8ArrayType:
		return ReadInt8Array(reader)
	case Int32ArrayType:
		return ReadInt32Array(reader)
	case ArrayType:
		return ReadArray(reader)
	case StringArrayType:
		return ReadStringArray(reader)
	case DictionaryType:
		return ReadDictionary(reader)
	case HashTableType:
		return ReadHashTable(reader)
	case NilType, UnknownType:
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported type: 0x%02x", ttype)
	}
}
