package v16

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"iter"
	"log"
	"math"
	"michelprogram/photon-parser/internal/types"
)

// Type represents a Photon Protocol16 type code.
// Each type code indicates how the following bytes should be interpreted.
type ParameterType uint8

// Photon Protocol16 type codes.
// These constants define the binary type codes used in Photon's serialization format.
const (
	UnknownType           ParameterType = 0x00 // Unknown or unsupported type
	NilType               ParameterType = 0x2a // Null/nil value
	DictionaryType        ParameterType = 0x44 // Dictionary with fixed key/value types
	StringArrayType       ParameterType = 0x61 // Array of strings
	Int8Type              ParameterType = 0x62 // 8-bit signed integer
	Custom                ParameterType = 0x63 // Custom serialized object
	DoubleType            ParameterType = 0x64 // Alias for Float64Type
	EventDateType         ParameterType = 0x65 // Event date/time
	Float32Type           ParameterType = 0x66 // 32-bit floating point
	Float64Type           ParameterType = 0x67 // 64-bit floating point
	HashTableType         ParameterType = 0x68 // Hashtable with mixed key/value types
	Int32Type             ParameterType = 0x69 // 32-bit signed integer
	Int16Type             ParameterType = 0x6b // 16-bit signed integer
	Int64Type             ParameterType = 0x6c // 64-bit signed integer
	Int32ArrayType        ParameterType = 0x6e // Array of 32-bit integers
	BooleanType           ParameterType = 0x6f // Boolean (0x00=false, 0x01=true)
	OperationResponseType ParameterType = 0x70 // Operation response message
	OperationRequestType  ParameterType = 0x71 // Operation request message
	StringType            ParameterType = 0x73 // UTF-8 string with uint16 length prefix
	Int8ArrayType         ParameterType = 0x78 // Array of 8-bit integers
	ArrayType             ParameterType = 0x79 // Generic typed array
	ObjectArrayType       ParameterType = 0x7a // Array of serialized objects
)

// Header represents the parameter header containing the parameter ID and type.
// This appears at the beginning of each serialized parameter.
type Header struct {
	ID   uint8         `json:"id"`   // Parameter identifier (application-specific)
	Type ParameterType `json:"type"` // Protocol16 type code indicating how to decode the value
}

// Parameters represents a complete Photon Protocol parameter with its header and decoded value.
// The Value field contains the decoded data according to the Type specified in the Header.
type Parameter struct {
	Header `json:"header"`
	Value  `json:"value"`
}

var _ types.ParameterView = (*Parameter)(nil)

type Value struct {
	Kind    ParameterType `json:"kind"`
	KeyType ParameterType `json:"key_type"`
	ValType ParameterType `json:"val_type"`
	_pad    [5]byte       `json:"-"`
	Num     uint64        `json:"num"`
	Str     string        `json:"str,omitempty"`
	Blob    []byte        `json:"blob,omitempty"`
}

func (p Parameter) ID() uint8 {
	return p.Header.ID
}

func (p Parameter) String() string {
	param := struct {
		Parameter `json:"parameter"`
	}{
		Parameter: p,
	}
	b, err := json.MarshalIndent(param, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(b)
}

func (p Parameter) Float32s() iter.Seq2[int, float32] {
	n := int(p.Num)
	if n <= 0 || len(p.Blob) < n*4 {
		return nil
	}
	return func(yield func(int, float32) bool) {
		for i := 0; i < n; i++ {
			bits := binary.LittleEndian.Uint32(p.Blob[i*4 : (i+1)*4])
			if !yield(i, math.Float32frombits(bits)) {
				return
			}
		}
	}
}

func (p Parameter) Float32() float32 {
	return math.Float32frombits(uint32(p.Num))
}

func (p Parameter) MarshalJSON() ([]byte, error) {
	type Alias Parameter

	out := struct {
		Alias
		Decoded any `json:"decoded,omitempty"`
	}{
		Alias: Alias(p),
	}

	// Example: special behavior by parameter kind/type
	switch p.Kind {
	case 5:
		log.Println("test")
		out.Decoded = p.Float32()
	case 69:
		res := make([]float32, p.Num)
		for index, fl := range p.Float32s() {
			res[index] = fl
		}
		out.Decoded = res
	default:
		out.Decoded = p.Num
	}

	return json.Marshal(out)
}

func (p Parameter) Float32ArrayValue() iter.Seq2[int, float32] {
	return nil
}
func (p Parameter) Int32ArrayValue() iter.Seq2[int, int32] {
	return nil
}
func (p Parameter) Int64ArrayValue() iter.Seq2[int, int64] {
	return nil
}
func (p Parameter) ByteArrayValue() iter.Seq2[int, byte] {
	return nil
}
func (p Parameter) Int16ArrayValue() iter.Seq2[int, int16] {
	return nil
}
func (p Parameter) StringArrayValue() iter.Seq2[int, string] {
	return nil
}
func (p Parameter) StringValue() string {
	return ""
}
func (p Parameter) Float32Value() float32 {
	return 0
}
func (p Parameter) IntValue() int64 {
	return 0
}
func (p Parameter) BooleanValue() bool {
	return false
}
func (p Parameter) ArrayValue() iter.Seq2[int, any] {
	return nil
}
func (p Parameter) BooleanArrayValue() iter.Seq2[int, bool] {
	return nil
}