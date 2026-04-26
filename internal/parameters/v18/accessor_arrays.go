package v18

import (
	"encoding/binary"
	"iter"
	"math"
	"michelprogram/photon-parser/internal/reader"
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

func (p Parameter) Int8ArrayValue() iter.Seq2[int, int8] {
	n := int(p.Num)
	if n <= 0 || len(p.Blob) < n || p.Kind != ByteArrayType {
		return nil
	}
	return func(yield func(int, int8) bool) {
		for i := 0; i < n; i++ {
			if !yield(i, int8(p.Blob[i])) {
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

func (p Parameter) ArrayValue() iter.Seq2[int, Value] {
    if p.Kind != ArrayType || p.Num == 0 || len(p.Blob) == 0 {
        return nil
    }
    return func(yield func(int, Value) bool) {
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
            if !yield(i, v) {
                return
            }
        }
    }
}