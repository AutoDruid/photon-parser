package v18

import (
	"encoding/binary"
	"iter"
	"math"

	"github.com/AutoDruid/photon-parser/internal/reader"
)

func (p Parameter) Float32ArrayValue() iter.Seq2[int, float32] {
	n := int(p.Num)
	if n <= 0 || len(p.Blob) < n*4 || p.Kind != Float32ArrayType {
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

func (p Parameter) Int32ArrayValue() iter.Seq2[int, int32] {
	if p.Kind != CompressedIntArrayType || p.Num == 0 || len(p.Blob) == 0 {
		return nil
	}

	return func(yield func(int, int32) bool) {
		r := reader.NewReader(p.Blob)

		for i := 0; i < int(p.Num); i++ {
			output, err := r.ReadVarintInt32()
			if err != nil {
				return
			}

			if !yield(i, output) {
				return
			}
		}
	}
}

func (p Parameter) Int64ArrayValue() iter.Seq2[int, int64] {
	if p.Kind != CompressedLongArrayType || p.Num == 0 || len(p.Blob) == 0 {
		return nil
	}

	return func(yield func(int, int64) bool) {
		r := reader.NewReader(p.Blob)

		for i := 0; i < int(p.Num); i++ {
			output, err := r.ReadVarintInt64()
			if err != nil {
				return
			}

			if !yield(i, output) {
				return
			}
		}
	}
}

func (p Parameter) ByteArrayValue() iter.Seq2[int, byte] {
	n := int(p.Num)
	if n <= 0 || len(p.Blob) < n || p.Kind != ByteArrayType {
		return nil
	}
	return func(yield func(int, byte) bool) {
		for i := 0; i < n; i++ {
			if !yield(i, p.Blob[i]) {
				return
			}
		}
	}
}

func (p Parameter) Int16ArrayValue() iter.Seq2[int, int16] {
	n := int(p.Num)
	if n <= 0 || len(p.Blob) < n*2 || p.Kind != ShortArrayType {
		return nil
	}
	return func(yield func(int, int16) bool) {
		for i := 0; i < n; i++ {
			value := int16(binary.LittleEndian.Uint16(p.Blob[i*2 : i*2+2]))
			if !yield(i, value) {
				return
			}
		}
	}
}

func (p Parameter) StringArrayValue() iter.Seq2[int, string] {
	if p.Kind != StringArrayType || p.Num == 0 || len(p.Blob) == 0 {
		return nil
	}

	return func(yield func(int, string) bool) {
		r := reader.NewReader(p.Blob)

		for i := 0; i < int(p.Num); i++ {
			size, err := r.ReadVarintUInt32()
			if err != nil {
				return
			}

			s, err := r.ReadString(int(size))
			if err != nil {
				return
			}

			if !yield(i, s) {
				return
			}
		}
	}
}

func (p Parameter) ArrayValue() iter.Seq2[int, any] {
	if p.Kind != ArrayType || p.Num == 0 || len(p.Blob) == 0 {
		return nil
	}
	return func(yield func(int, any) bool) {
		r := reader.NewReader(p.Blob)
		for i := 0; i < int(p.Num); i++ {
			ttype, err := r.ReadByte()
			if err != nil {
				return
			}

			v, err := scanPayload(r, ParameterType(ttype))
			if err != nil {
				return
			}

			if !yield(i, decodeValue(v)) {
				return
			}
		}
	}
}

func (p Parameter) BooleanArrayValue() iter.Seq2[int, bool] {
	if p.Kind != BooleanArrayType || p.Num == 0 || len(p.Blob) == 0 {
		return nil
	}

	return func(yield func(int, bool) bool) {
		for i := 0; i < int(p.Num); i++ {
			out := (p.Blob[i/8] & (1 << uint(i%8))) != 0
			if !yield(i, out) {
				return
			}
		}
	}
}

func decodeValue(v Value) any {
	p := Parameter{Value: v}

	switch v.Kind {
	case StringType:
		v, _ := p.StringValue()
		return v
	case Float32Type:
		v, _ := p.Float32Value()
		return v
	case BooleanTrueType, BooleanFalseType, BooleanType:
		v, _ := p.BooleanValue()
		return v
	case Int8Type, Int8Positive, Int8Negative,
		Int16Type, Int16Positive, Int16Negative,
		CompressedInt32Type, CompressedInt64Type,
		Long8Positive, Long8Negative, Long16Positive, Long16Negative,
		IntZeroType, ShortZeroType, LongZeroType, ByteZeroType:
		v, _ := p.IntValue()
		return v
	case ByteArrayType:
		return collect(p.ByteArrayValue(), v.Num)
	case ShortArrayType:
		return collect(p.Int16ArrayValue(), v.Num)
	case CompressedIntArrayType:
		return collect(p.Int32ArrayValue(), v.Num)
	case CompressedLongArrayType:
		return collect(p.Int64ArrayValue(), v.Num)
	case Float32ArrayType:
		return collect(p.Float32ArrayValue(), v.Num)
	case StringArrayType:
		return collect(p.StringArrayValue(), v.Num)
	case BooleanArrayType:
		return collect(p.BooleanArrayValue(), v.Num)
	case ArrayType:
		return collect(p.ArrayValue(), v.Num)
	case NilType:
		return nil
	case DictionaryType:
		return collectDictionary(p.DictionaryValue(), v.Num)
	default:
		return v
	}
}

func (p Parameter) Int8ArrayValue() iter.Seq2[int, int8] {
	return nil
}

func collectDictionary(seq iter.Seq2[any, any], n uint64) map[any]any {
	if seq == nil {
		return nil
	}

	out := make(map[any]any, n)
	for key, value := range seq {
		out[key] = value
	}
	return out
}

func collect[T any](seq iter.Seq2[int, T], n uint64) []T {
	if seq == nil {
		return nil
	}

	out := make([]T, 0, n)
	for _, value := range seq {
		out = append(out, value)
	}
	return out
}
