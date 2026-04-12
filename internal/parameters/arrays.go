package parameters

import (
	"michelprogram/photon-parser/internal/reader"
)

// ReadInt8Array reads an array of 8-bit signed integers from the reader.
// Format: uint32 size followed by size int8 values.
// Returns an error if the array cannot be fully read.
func (p Parameters) readInt8Array(r *reader.Reader) ([]int8, error) {
	size, err := r.ReadUInt32()
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

// ReadInt32Array reads an array of 32-bit signed integers from the reader.
// Format: uint32 size followed by size int32 values (each in big-endian).
// Returns an error if the array cannot be fully read.
func (p Parameters) readInt32Array(r *reader.Reader) ([]int32, error) {
	size, err := r.ReadUInt32()
	if err != nil {
		return nil, err
	}

	val := make([]int32, size)

	for i := uint32(0); i < size; i++ {
		input, err := r.ReadInt32()
		if err != nil {
			return nil, err
		}
		val[i] = input
	}
	return val, nil
}

// ReadStringArray reads an array of strings from the reader.
// Format: uint32 size followed by size Protocol16 strings (each with uint16 length prefix).
// Returns an error if the array cannot be fully read.
func (p Parameters) readStringArray(r *reader.Reader) ([]string, error) {
	size, err := r.ReadUInt32()
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

// ReadArray reads a generic typed array from the reader.
// Format: uint16 size, Type byte, then size elements of that type.
// All elements in the array have the same type.
// Returns a slice of any containing the decoded elements.
//
// Example wire format for []int32{100, 200}:
//
//	0x00 0x02  // size = 2 (uint16)
//	0x69                 // type = Int32Type
//	0x00 0x00 0x00 0x64  // 100
//	0x00 0x00 0x00 0xC8  // 200
func (p Parameters) readArray(r *reader.Reader) ([]any, error) {

	size, err := r.ReadUInt16()
	if err != nil {
		return nil, err
	}

	ttype, err := r.ReadUInt8()
	if err != nil {
		return nil, err
	}

	val := make([]any, size)

	for i := uint16(0); i < size; i++ {
		input, err := p.decode(r, Type(ttype))
		if err != nil {
			return nil, err
		}
		val[i] = input
	}
	return val, nil
}
