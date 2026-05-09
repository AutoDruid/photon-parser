package v16

import (
	"fmt"

	"github.com/AutoDruid/photon-parser/internal/context"
	"github.com/AutoDruid/photon-parser/internal/hooks"
	"github.com/AutoDruid/photon-parser/internal/reader"
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
		return err
	}

	out.Header = header
	out.Value = value

	p.emit(reader, hooks)

	return nil
}

func (p Parameter) emit(reader *reader.Reader, hooks *hooks.Hooks[Parameter]) {
	if hooks == nil {
		return
	}

	if hooks.SyncHooks.OnParameter != nil {
		hooks.SyncHooks.OnParameter(p)
	}

	if hooks.AsyncHooks.OnParameter == nil {
		return
	}

	select {
	case hooks.AsyncHooks.OnParameter <- p:
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
func scanPayload(reader *reader.Reader, t ParameterType) (Value, error) {
	var err error
	res := Value{Kind: t}

	switch t {
	case StringType:
		err = scanString(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case Float32Type:
		err = scanFloat32(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case Float64Type:
		err = scanFloat64(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case Int8Type:
		err = scanInt8(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case Int16Type:
		err = scanInt16(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case Int32Type:
		err = scanInt32(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case Int64Type:
		err = scanInt64(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case BooleanType:
		err = scanBoolean(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case Int8ArrayType:
		err = scanInt8Array(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case Int32ArrayType:
		err = scanInt32Array(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case StringArrayType:
		err = scanStringArray(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case ArrayType:
		err = scanArray(reader, &res)
		if err != nil {
			return Value{}, err
		}
	case NilType, UnknownType:
		break
	case DictionaryType:
		err = scanDictionary(reader, &res)
		if err != nil {
			return Value{}, err
		}
	default:
		return Value{}, fmt.Errorf("unsupported type: %d", t)
	}

	return res, nil
}
