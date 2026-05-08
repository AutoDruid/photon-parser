package v16

import (
	"encoding/binary"
	"fmt"
	"iter"
	"michelprogram/photon-parser/internal/reader"
)

func (p Parameter) ByteArrayValue() iter.Seq2[int, byte] {
	return func(yield func(int, byte) bool) {
		n := int(p.Num)

		if n <= 0 || len(p.Blob) < n || p.Kind != Int8ArrayType {
			return
		}

		for i := 0; i < n; i++ {
			if !yield(i, p.Blob[i]) {
				return
			}
		}
	}
}

func (p Parameter) Int8ArrayValue() iter.Seq2[int, int8] {
	return func(yield func(int, int8) bool) {
		if p.Kind != Int8ArrayType || p.Num == 0 || len(p.Blob) == 0 {
			return
		}

		r := reader.NewReader(p.Blob)

		for i := 0; i < int(p.Num); i++ {
			output, err := r.ReadInt8()
			if err != nil {
				return
			}

			if !yield(i, output) {
				return
			}
		}
	}
}

func (p Parameter) Int32ArrayValue() iter.Seq2[int, int32] {
	return func(yield func(int, int32) bool) {
		if p.Kind != Int32ArrayType || p.Num == 0 || len(p.Blob) == 0 {
			return
		}

		r := reader.NewReader(p.Blob)

		for i := 0; i < int(p.Num); i++ {
			output, err := r.ReadInt32(binary.BigEndian)
			if err != nil {
				return
			}

			if !yield(i, output) {
				return
			}
		}
	}
}

func (p Parameter) StringArrayValue() iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		if p.Kind != StringArrayType || p.Num == 0 || len(p.Blob) == 0 {
			return
		}

		r := reader.NewReader(p.Blob)

		for i := 0; i < int(p.Num); i++ {
			size, err := r.ReadUInt16(binary.BigEndian)
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
	return func(yield func(int, any) bool) {
		if p.Kind != ArrayType || p.Num == 0 || len(p.Blob) == 0 {
			return
		}

		r := reader.NewReader(p.Blob)

		for i := 0; i < int(p.Num); i++ {

			v, err := scanPayload(r, p.KeyType)
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
	return func(yield func(int, bool) bool) {
		if p.Kind != ArrayType && p.KeyType != BooleanType || p.Num == 0 || len(p.Blob) == 0 {
			return
		}

		r := reader.NewReader(p.Blob)
		for i := 0; i < int(p.Num); i++ {
			b, err := r.ReadBoolean()
			if err != nil {
				return
			}

			if !yield(i, b) {
				return
			}
		}
	}
}

func (p Parameter) Float32ArrayValue() iter.Seq2[int, float32] {
	return nil
}

func (p Parameter) Int64ArrayValue() iter.Seq2[int, int64] {
	return nil
}

func (p Parameter) Int16ArrayValue() iter.Seq2[int, int16] {
	return nil
}

func decodeValue(v Value) any {
	p := Parameter{Value: v}

	switch v.Kind {
	case StringType:
		v, _ := p.StringValue()
		return v
	case Int8Type, Int16Type, Int32Type, Int64Type:
		v, _ := p.IntValue()
		return v
	case BooleanType:
		v, _ := p.BooleanValue()
		return v
	case ArrayType:
		return collect(p.ArrayValue(), p.Num)
	default:
		return fmt.Errorf("unknown type: %d\n", v.Kind)
	}
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
