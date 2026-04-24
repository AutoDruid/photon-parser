package v18

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"michelprogram/photon-parser/internal/context"
	"michelprogram/photon-parser/internal/hooks"
	"michelprogram/photon-parser/internal/reader"
	"michelprogram/photon-parser/internal/types"
)

type Parameter struct {
	types.Parameter
}

var _ context.ParameterParser = (*Parameter)(nil)

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

	value, err := p.decode(reader, ParameterType(header.Type))

	if err != nil {
		log.Println("err on parameter type", header.Type, err)
		return err
	}

	out.ParameterHeader = header
	out.Value = value

	//p.emit(reader, hooks, out)

	return nil
}

func (p Parameter) emit(reader *reader.Reader, hooks *hooks.Hooks, out *types.Parameter) {
	if hooks == nil {
		return
	}

	if hooks.SyncHooks.OnParameter != nil {
		hooks.SyncHooks.OnParameter(*out)
	}

	if hooks.AsyncHooks.OnParameter == nil {
		return
	}

	select {
	case hooks.AsyncHooks.OnParameter <- *out:
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

func (p Parameter) decode(reader *reader.Reader, t ParameterType) (types.Value, error) {
	res := types.Value{Kind: types.ParameterType(t)}

	switch t {
	case Int8Positive:
		b, err := reader.ReadByte()
		if err != nil {
			return types.Value{}, err
		}
		res.Num = uint64(uint(int32(b)))
	case Int8Negative:
		b, err := reader.ReadByte()
		if err != nil {
			return types.Value{}, err
		}
		res.Num = uint64(uint32(-int32(b)))
	case Int16Type:
		value, err := reader.ReadInt16(binary.LittleEndian)
		if err != nil {
			return types.Value{}, err
		}
		res.Num = uint64(uint16(value))
	case Int16Positive:
		value, err := reader.ReadUInt16(binary.LittleEndian)
		if err != nil {
			return types.Value{}, err
		}
		res.Num = uint64(uint32(int32(value)))
	case Int16Negative:
		value, err := reader.ReadUInt16(binary.LittleEndian)
		if err != nil {
			return types.Value{}, err
		}
		res.Num = uint64(uint32(-int32(value)))
	case Long8Positive:
		value, err := reader.ReadByte()
		if err != nil {
			return types.Value{}, err
		}
		res.Num = uint64(value)
	case Long8Negative:
		value, err := reader.ReadByte()
		if err != nil {
			return types.Value{}, err
		}
		res.Num = uint64(-int64(value))
	case Long16Positive:
		value, err := reader.ReadUInt16(binary.LittleEndian)
		if err != nil {
			return types.Value{}, err
		}
		res.Num = uint64(value)
	case Long16Negative:
		value, err := reader.ReadUInt16(binary.LittleEndian)
		if err != nil {
			return types.Value{}, err
		}
		res.Num = uint64(-int64(value))
	case StringType:
		value, err := p.readString(reader)
		if err != nil {
			return types.Value{}, err
		}
		res.Str = value
	case CompressedInt32Type:
		value, err := reader.ReadVarintInt32()
		if err != nil {
			return types.Value{}, err
		}
		res.Num = uint64(uint32(value))
	case CompressedInt64Type:
		value, err := reader.ReadVarintInt64()
		if err != nil {
			return types.Value{}, err
		}
		res.Num = uint64(value)
	case Float32ArrayType:
		count, err := reader.ReadVarintUInt32()
		if err != nil {
			return types.Value{}, err
		}
		blob, err := p.readBlob(reader, int(count)*4)
		if err != nil {
			return types.Value{}, err
		}
		res.Blob = blob
		res.Num = uint64(count)
	case Float32Type:
		value, err := reader.ReadFloat32(binary.LittleEndian)
		if err != nil {
			return types.Value{}, err
		}
		res.Num = uint64(math.Float32bits(value))
	case Int8Type:
		value, err := reader.ReadInt8()
		if err != nil {
			return types.Value{}, err
		}
		res.Num = uint64(uint8(value))
	case BooleanTrueType:
		res.Num = 1
	case BooleanFalseType:
		res.Num = 0
	case IntZeroType, ShortZeroType, ByteZeroType:
		break
	case ArrayType:
		count, err := reader.ReadVarintUInt32()
		if err != nil {
			return types.Value{}, err
		}
		start := reader.Cursor
		for i := uint32(0); i < count; i++ {
			ttype, err := reader.ReadByte()
			if err != nil {
				return types.Value{}, err
			}
			if _, err := p.decode(reader, ParameterType(ttype)); err != nil {
				return types.Value{}, err
			}
		}
		res.Blob = reader.Buffer[start:reader.Cursor]
		res.Num = uint64(count)
	case ShortArrayType:
		count, err := reader.ReadVarintUInt32()
		if err != nil {
			return types.Value{}, err
		}
		blob, err := p.readBlob(reader, int(count)*2)
		if err != nil {
			return types.Value{}, err
		}
		res.Blob = blob
		res.Num = uint64(count)
	case ByteArrayType:
		count, err := reader.ReadVarintUInt32()
		if err != nil {
			return types.Value{}, err
		}
		blob, err := p.readBlob(reader, int(count))
		if err != nil {
			return types.Value{}, err
		}
		res.Blob = blob
		res.Num = uint64(count)
	case BooleanArrayType:
		count, err := reader.ReadVarintUInt32()
		if err != nil {
			return types.Value{}, err
		}
		packedBytes := (int(count) + 7) / 8
		blob, err := p.readBlob(reader, packedBytes)
		if err != nil {
			return types.Value{}, err
		}
		res.Blob = blob
		res.Num = uint64(count)
	case StringArrayType:
		count, err := reader.ReadVarintUInt32()
		if err != nil {
			return types.Value{}, err
		}
		start := reader.Cursor
		for i := uint32(0); i < count; i++ {
			if _, err := p.readString(reader); err != nil {
				return types.Value{}, err
			}
		}
		res.Blob = reader.Buffer[start:reader.Cursor]
		res.Num = uint64(count)
	case DictionaryType:
		keyType, err := reader.ReadUInt8()
		if err != nil {
			return types.Value{}, err
		}
		valueType, err := reader.ReadUInt8()
		if err != nil {
			return types.Value{}, err
		}
		count, err := reader.ReadVarintUInt32()
		if err != nil {
			return types.Value{}, err
		}
		start := reader.Cursor
		for i := uint32(0); i < count; i++ {
			if _, err := p.decode(reader, ParameterType(keyType)); err != nil {
				return types.Value{}, err
			}
			if _, err := p.decode(reader, ParameterType(valueType)); err != nil {
				return types.Value{}, err
			}
		}
		res.Blob = reader.Buffer[start:reader.Cursor]
		res.Num = uint64(count)
		res.KeyType = types.ParameterType(keyType)
		res.ValType = types.ParameterType(valueType)
	case CompressedIntArrayType:
		count, err := reader.ReadVarintUInt32()
		if err != nil {
			return types.Value{}, err
		}
		blob := make([]byte, int(count)*4)
		for i := 0; i < int(count); i++ {
			n, err := reader.ReadVarintInt32()
			if err != nil {
				return types.Value{}, err
			}
			binary.BigEndian.PutUint32(blob[i*4:], uint32(n))
		}
		res.Blob = blob
		res.Num = uint64(count)
	case CompressedLongArrayType:
		count, err := reader.ReadVarintUInt32()
		if err != nil {
			return types.Value{}, err
		}
		blob := make([]byte, int(count)*8)
		for i := 0; i < int(count); i++ {
			n, err := reader.ReadVarintInt64()
			if err != nil {
				return types.Value{}, err
			}
			binary.BigEndian.PutUint64(blob[i*8:], uint64(n))
		}
		res.Blob = blob
		res.Num = uint64(count)
	case NilType, UnknownType:
		break
	default:
		return types.Value{}, fmt.Errorf("unsupported type: %d", t)
	}
	return res, nil
}

func (p Parameter) readBlob(r *reader.Reader, n int) ([]byte, error) {
	raw, err := r.ReadBytes(n)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

// String returns a human-readable representation of the parameter.
// Format: "ID: <id>\nType: <type>\nValue: <value>\n"
func (p Parameter) String() string {
	return fmt.Sprintf("ID: %d\nType: %d\nValue: %v\n", p.ID, p.Type, p.Value)
}
