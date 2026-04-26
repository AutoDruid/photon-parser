package v18

import (
	"fmt"
	"log"
	"michelprogram/photon-parser/internal/context"
	"michelprogram/photon-parser/internal/hooks"
	"michelprogram/photon-parser/internal/reader"
)

var _ context.ParameterParser[Parameter] = (*Parameter)(nil)

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
func (p *Parameter) Parse(reader *reader.Reader, out *Parameter, hooks *hooks.Hooks[Parameter]) error {

	header, err := p.parseHeader(reader)
	if err != nil {
		return err
	}

	value, err := scanPayload(reader, header.Type)

	if err != nil {
		log.Println("err on parameter type", header.Type, err)
		return err
	}

	out.Header = header
	out.Value = value

	//p.emit(reader, hooks, out)

	return nil
}

func (p Parameter) emit(reader *reader.Reader, hooks *hooks.Hooks[Parameter], out *Parameter) {
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

func (p *Parameter) parseHeader(r *reader.Reader) (Header, error) {
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

	header.Type = ParameterType(b)

	return header, nil
}

func scanPayload(reader *reader.Reader, t ParameterType) (Value, error) {
	var err error
	var res Value = Value{Kind: t}

	switch t {
	case Int8Type:
		err = scanInt8(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case Int8Positive:
		err = scanInt8Positive(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case Int8Negative:
		err = scanInt8Negative(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case Int16Type:
		err = scanInt16Type(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case Int16Positive:
		err = scanInt16Positive(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case Int16Negative:
		err = scanInt16Negative(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case Long8Positive:
		err = scanLong8Positive(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case Long8Negative:
		err = scanLong8Negative(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case Long16Positive:
		err = scanLong16Positive(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case Long16Negative:
		err = scanLong16Negative(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case StringType:
		err = scanString(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case CompressedInt32Type:
		err = scanCompressedInt32(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case CompressedInt64Type:
		err = scanCompressedInt64(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case Float32ArrayType:
		err = scanFloat32Array(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case Float32Type:
		err = scanFloat32(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case BooleanTrueType:
		res.Num = 1
	case BooleanFalseType:
		res.Num = 0
	case IntZeroType, ShortZeroType, ByteZeroType:
		break
	case ArrayType:
		err = scanArray(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case ShortArrayType:
		err = scanShortArray(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case ByteArrayType:
		err = scanByteArray(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case BooleanArrayType:
		err = scanBooleanArray(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case StringArrayType:
		err = scanStringArray(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case DictionaryType:
		err = scanDictionary(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case CompressedIntArrayType:
		err = scanCompressedIntArray(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case CompressedLongArrayType:
		err = scanCompressedLongArray(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case NilType, UnknownType:
		break
	default:
		return Value{}, fmt.Errorf("unsupported type: %d", t)
	}
	return res, nil
}
