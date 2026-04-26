package v18

import (
	"encoding/json"
	"michelprogram/photon-parser/internal/types"
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
	Blob    []byte        `json:"blob,omitempty"`
}

func (p Parameter) ID() uint8 {
	return p.Header.ID
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
		out.Decoded = p.Float32Value()
	case 69:
		res := make([]float32, p.Num)
		for index, fl := range p.Float32ArrayValue() {
			res[index] = fl
		}
		out.Decoded = res
	default:
		out.Decoded = p.Num
	}

	return json.Marshal(out)
}
