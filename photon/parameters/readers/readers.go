package readers

import (
	"bytes"
)

func Decode(reader *bytes.Reader, ttype Type) (any, error){
	switch ttype {
		default:
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
	}

	return "",nil
}
