package v18

import (
	"encoding/binary"
	"michelprogram/photon-parser/internal/reader"
)

func (p Parameter) readFloatArray(r *reader.Reader) ([]float32, error) {
	size, err := r.ReadVarintUInt32()
	if err != nil {
		return nil, err
	}

	val := make([]float32, size)

	for i := uint32(0); i < size; i++ {
		input, err := r.ReadFloat32(binary.BigEndian)
		if err != nil {
			return nil, err
		}
		val[i] = input
	}
	return val, nil
}

func (p Parameter) readInt8Array(r *reader.Reader) ([]int8, error) {
	size, err := r.ReadVarintUInt32()
	if err != nil {
		return nil, err
	}

	val := make([]int8, size)

	for i := uint32(0); i < size; i++ {
		input, err := r.ReadInt8()
		if err != nil {
			return nil, err
		}
		val[i] = input
	}
	return val, nil
}

func (p Parameter) readInt16Array(r *reader.Reader) ([]int16, error) {
	size, err := r.ReadVarintUInt32()
	if err != nil {
		return nil, err
	}

	val := make([]int16, size)

	for i := uint32(0); i < size; i++ {
		input, err := r.ReadInt16(binary.LittleEndian)
		if err != nil {
			return nil, err
		}
		val[i] = input
	}
	return val, nil
}

func (p Parameter) readStringArray(r *reader.Reader) ([]string, error) {
	size, err := r.ReadVarintUInt32()
	if err != nil {
		return nil, err
	}

	val := make([]string, size)

	for i := uint32(0); i < size; i++ {
		input, err := p.readString(r)
		if err != nil {
			return nil, err
		}
		val[i] = input
	}
	return val, nil
}

func (p Parameter) readCompressedInt32Array(r *reader.Reader) ([]int32, error) {

	size, err := r.ReadVarintUInt32()
	if err != nil {
		return nil, err
	}

	val := make([]int32, size)

	for i := uint32(0); i < size; i++ {
		input, err := r.ReadVarintInt32()
		if err != nil {
			return nil, err
		}
		val[i] = input
	}
	return val, nil
}

func (p Parameter) readCompressedInt64Array(r *reader.Reader) ([]int64, error) {
	size, err := r.ReadVarintUInt32()
	if err != nil {
		return nil, err
	}

	val := make([]int64, size)

	for i := uint32(0); i < size; i++ {
		input, err := r.ReadVarintInt64()
		if err != nil {
			return nil, err
		}
		val[i] = input
	}
	return val, nil
}

func (p Parameter) readArray(r *reader.Reader) ([]interface{}, error) {
	size, err := r.ReadVarintUInt32()
	if err != nil {
		return nil, err
	}

	val := make([]interface{}, size)
	for i := uint32(0); i < size; i++ {
		ttype, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		input, err := p.decode(r, ParameterType(ttype))
		if err != nil {
			return nil, err
		}
		val[i] = input
	}
	return val, nil
}

// ReadPackedBooleanArray reads size booleans packed 8-per-byte (LSB-first).
// It advances r.Cursor by ceil(size/8) bytes.
func (p Parameter) readBooleanArray(r *reader.Reader) ([]bool, error) {
	size, err := r.ReadVarintUInt32()
	if err != nil {
		return nil, err
	}

	packedBytes := (int(size) + 7) / 8
	packed, err := r.ReadBytes(packedBytes)
	if err != nil {
		return nil, err
	}

	out := make([]bool, size)
	for i := uint32(0); i < size; i++ {
		out[i] = (packed[i/8] & (1 << uint(i%8))) != 0
	}
	return out, nil
}
