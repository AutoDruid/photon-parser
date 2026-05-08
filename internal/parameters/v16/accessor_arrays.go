package v16

import (
	"encoding/binary"
	"iter"
	"log"
	"michelprogram/photon-parser/internal/reader"
)

func (p Parameter) ByteArrayValue() iter.Seq2[int, byte] {
	n := int(p.Num)
	if n <= 0 || len(p.Blob) < n || p.Kind != Int8ArrayType {
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

		for i := 0; i < int(p.Num); i++ {
			r := reader.NewReader(p.Blob)

			size, err := r.ReadUInt16(binary.BigEndian)
			if err != nil {
				return
			}

			s, err := r.ReadString(int(size))
			if err != nil {
				return
			}

			log.Println("Readed", s, size)

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
		ttype, err := r.ReadUInt8()
		if err != nil {
			return
		}
		for i := 0; i < int(p.Num); i++ {

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
	if p.Kind != ArrayType || p.Num == 0 || len(p.Blob) == 0 {
		return nil
	}

	r := reader.NewReader(p.Blob)
	ttype, err := r.ReadUInt8()
	if err != nil {
		return nil
	}

	if ParameterType(ttype) != BooleanType {
		return nil
	}

	return func(yield func(int, bool) bool) {
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

func decodeValue(v Value) any {
	p := Parameter{Value: v}

	switch v.Kind {
	case StringType:
		v, _ := p.StringValue()
		return v
	}
	return v
}
