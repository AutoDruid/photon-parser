package v18

import (
	"michelprogram/photon-parser/internal/reader"
)

func (p Parameter) readFloatArray(r *reader.Reader) ([]float32, error) {
	size, err := r.ReadVarintUInt32()
	if err != nil {
		return nil, err
	}

	val := make([]float32, size)

	for i := uint32(0); i < size; i++ {
		input, err := r.ReadFloat32()
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
		input, err := r.ReadInt16LittleEndian()
		if err != nil {
			return nil, err
		}
		val[i] = input
	}
	return val, nil
}
