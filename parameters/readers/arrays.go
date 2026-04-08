package readers

import (
	"michelprogram/photon-parser/parser"

	"golang.org/x/exp/constraints"
)

// readPrimitiveArray is a generic helper that reads an array of primitive numeric types.
// Format: uint32 size (big-endian) followed by size elements of type T.
// This is used internally by ReadInt8Array, ReadInt32Array, etc.
func readPrimitiveArray[T constraints.Integer | constraints.Float](reader *parser.Reader) ([]T, error) {
	size, err := parser.ReadPrimitive[uint32](reader)

	if err != nil {
		return nil, err
	}

	val := make([]T, size)

	for i := uint32(0); i < size; i++ {
		input, err := parser.ReadPrimitive[T](reader)
		if err != nil {
			return nil, err
		}
		val[i] = input
	}
	return val, nil
}

// ReadInt8Array reads an array of 8-bit signed integers from the reader.
// Format: uint32 size followed by size int8 values.
// Returns an error if the array cannot be fully read.
func ReadInt8Array(reader *parser.Reader) ([]int8, error) {
	return readPrimitiveArray[int8](reader)
}

// ReadInt32Array reads an array of 32-bit signed integers from the reader.
// Format: uint32 size followed by size int32 values (each in big-endian).
// Returns an error if the array cannot be fully read.
func ReadInt32Array(reader *parser.Reader) ([]int32, error) {
	return readPrimitiveArray[int32](reader)
}

// ReadStringArray reads an array of strings from the reader.
// Format: uint32 size followed by size Protocol16 strings (each with uint16 length prefix).
// Returns an error if the array cannot be fully read.
func ReadStringArray(reader *parser.Reader) ([]string, error) {
	size, err := parser.ReadPrimitive[uint32](reader)

	if err != nil {
		return nil, err
	}

	val := make([]string, size)

	for i := uint32(0); i < size; i++ {
		input, err := ReadString(reader)
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
func ReadArray(reader *parser.Reader) ([]any, error) {

	size, err := parser.ReadPrimitive[uint16](reader)

	if err != nil {
		return nil, err
	}

	ttype, err := parser.ReadPrimitive[Type](reader)

	if err != nil {
		return nil, err
	}

	val := make([]any, size)

	for i := uint16(0); i < size; i++ {
		input, err := Decode(reader, ttype)
		if err != nil {
			return nil, err
		}
		val[i] = input
	}
	return val, nil
}
