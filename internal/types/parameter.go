package types

import "iter"

type ParameterView interface {
	ID() uint8

	Float32ArrayValue() iter.Seq2[int, float32]
	Int32ArrayValue() iter.Seq2[int, int32]
	Int64ArrayValue() iter.Seq2[int, int64]
	ByteArrayValue() iter.Seq2[int, byte]
	Int16ArrayValue() iter.Seq2[int, int16]
	StringArrayValue() iter.Seq2[int, string]
	BooleanArrayValue() iter.Seq2[int, bool]
	ArrayValue() iter.Seq2[int, any]

	BooleanValue() (bool, bool)
	StringValue() (string, bool)
	Float32Value() (float32, bool)
	IntValue() (int64, bool)

	
	MarshalJSON() ([]byte, error)
}
