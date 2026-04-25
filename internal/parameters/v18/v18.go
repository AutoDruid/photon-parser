package v18

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"iter"
	"log"
	"math"
	"michelprogram/photon-parser/internal/types"
)

type ParameterType uint8

const (
	UnknownType         ParameterType = 0 //Tested
	BooleanType         ParameterType = 2
	Int8Type            ParameterType = 3
	Int16Type           ParameterType = 4 //Tested
	Float32Type         ParameterType = 5
	Float64Type         ParameterType = 6
	StringType          ParameterType = 7 //Tested
	NilType             ParameterType = 8
	CompressedInt32Type ParameterType = 9  //Tested
	CompressedInt64Type ParameterType = 10 //Tested

	Int8Positive  ParameterType = 11 // 1 byte unsigned, cast to +int32
	Int8Negative  ParameterType = 12 // 1 byte unsigned, cast to -int32
	Int16Positive ParameterType = 13 // 2 bytes unsigned, cast to +int32 Tested
	Int16Negative ParameterType = 14 // 2 bytes unsigned, cast to -int32

	Long8Positive  ParameterType = 15 // 1 byte unsigned, cast to +int64
	Long8Negative  ParameterType = 16 // 1 byte unsigned, cast to -int64
	Long16Positive ParameterType = 17 // 2 bytes unsigned, cast to +int64
	Long16Negative ParameterType = 18 // 2 bytes unsigned, cast to -int64

	CustomType     ParameterType = 19
	CustomTypeSlim ParameterType = 0x80

	// Complex types
	DictionaryType        ParameterType = 20
	HashtableType         ParameterType = 21
	ObjectArrayType       ParameterType = 23
	OperationRequestType  ParameterType = 24
	OperationResponseType ParameterType = 25
	EventDataType         ParameterType = 26

	// Zero shorthands — no payload bytes, type code is the entire value
	BooleanFalseType ParameterType = 27
	BooleanTrueType  ParameterType = 28
	ShortZeroType    ParameterType = 29
	IntZeroType      ParameterType = 30
	LongZeroType     ParameterType = 31
	FloatZeroType    ParameterType = 32
	DoubleZeroType   ParameterType = 33
	ByteZeroType     ParameterType = 34 //Tested

	// Array container — element count + element type follows
	ArrayType ParameterType = 0x40

	// Typed arrays — element type is baked into the type code (elemType | 0x40)
	BooleanArrayType        ParameterType = BooleanType | ArrayType         // 0x42
	ByteArrayType           ParameterType = Int8Type | ArrayType            // 0x43
	ShortArrayType          ParameterType = Int16Type | ArrayType           // 0x44 Tested
	Float32ArrayType        ParameterType = Float32Type | ArrayType         // 0x45 Tested
	Float64ArrayType        ParameterType = Float64Type | ArrayType         // 0x46
	StringArrayType         ParameterType = StringType | ArrayType          // 0x47
	CompressedIntArrayType  ParameterType = CompressedInt32Type | ArrayType // 0x49
	CompressedLongArrayType ParameterType = CompressedInt64Type | ArrayType // 0x4A
	CustomTypeArrayType     ParameterType = CustomType | ArrayType          // 0x53
	DictionaryArrayType     ParameterType = DictionaryType | ArrayType      // 0x54
	HashtableArrayType      ParameterType = HashtableType | ArrayType       // 0x55
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

var _ types.VersionedParameter = (*Parameter)(nil)

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
