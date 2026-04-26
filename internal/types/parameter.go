package types

import "iter"

type ParameterView interface {
	ID() uint8
	Float32ArrayValue() iter.Seq2[int, float32]
	Int32ArrayValue() iter.Seq2[int, int32]
	Int64ArrayValue() iter.Seq2[int, int64]
	Int8ArrayValue() iter.Seq2[int, int8]
	Int16ArrayValue() iter.Seq2[int, int16]
	StringArrayValue() iter.Seq2[int, string]
	StringValue() string
	Float32Value() float32
	IntValue() int64
	MarshalJSON() ([]byte, error)
}
