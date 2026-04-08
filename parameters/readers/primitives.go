package readers

import (
	"fmt"
	"io"
	"michelprogram/photon-parser/parser"
	"sync"
)

// stringBufPool holds reusable byte slices used as temporaries in ReadString.
// ReadString converts the slice to a string immediately (copying the bytes),
// so the pooled buffer is safe to reuse right after string conversion.
var stringBufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 64)
		return &b
	},
}

// ReadInt8 reads an 8-bit signed integer from the reader.
// Returns an error if fewer than 1 byte is available.
func ReadInt8(reader *parser.Reader) (int8, error) {
	return parser.ReadPrimitive[int8](reader)
}

// ReadInt16 reads a 16-bit signed integer from the reader in big-endian format.
// Returns an error if fewer than 2 bytes are available.
func ReadInt16(reader *parser.Reader) (int16, error) {
	return parser.ReadPrimitive[int16](reader)
}

// ReadInt32 reads a 32-bit signed integer from the reader in big-endian format.
// Returns an error if fewer than 4 bytes are available.
func ReadInt32(reader *parser.Reader) (int32, error) {
	return parser.ReadPrimitive[int32](reader)
}

// ReadInt64 reads a 64-bit signed integer from the reader in big-endian format.
// Returns an error if fewer than 8 bytes are available.
func ReadInt64(reader *parser.Reader) (int64, error) {
	return parser.ReadPrimitive[int64](reader)
}

// ReadFloat32 reads a 32-bit floating point number from the reader in big-endian format.
// Returns an error if fewer than 4 bytes are available.
func ReadFloat32(reader *parser.Reader) (float32, error) {
	return parser.ReadPrimitive[float32](reader)
}

// ReadFloat64 reads a 64-bit floating point number from the reader in big-endian format.
// Returns an error if fewer than 8 bytes are available.
func ReadFloat64(reader *parser.Reader) (float64, error) {
	return parser.ReadPrimitive[float64](reader)
}

// ReadString reads a Photon Protocol16 string from the reader.
// Format: uint16 length (big-endian) followed by UTF-8 bytes.
// Returns an empty string if length is 0.
// Returns an error if the declared length cannot be fully read.
//
// Example wire format for "hello":
//
//	0x00 0x05 'h' 'e' 'l' 'l' 'o'
func ReadString(reader *parser.Reader) (string, error) {
	size, err := ReadInt16(reader)
	if err != nil {
		return "", err
	}

	if size == 0 {
		return "", nil
	}

	bufp := stringBufPool.Get().(*[]byte)
	buf := *bufp
	if cap(buf) < int(size) {
		buf = make([]byte, size)
	}
	buf = buf[:size]

	if _, err := io.ReadFull(reader, buf); err != nil {
		stringBufPool.Put(bufp)
		return "", fmt.Errorf("failed to read string: %w", err)
	}
	s := string(buf)
	*bufp = buf
	stringBufPool.Put(bufp)
	return s, nil
}

// ReadBoolean reads a boolean value from the reader.
// Format: single byte (0x00 = false, 0x01 = true).
// Returns an error if the value is not 0x00 or 0x01, or if no byte is available.
func ReadBoolean(readers *parser.Reader) (bool, error) {
	value, err := parser.ReadPrimitive[uint8](readers)

	if err != nil {
		return false, err
	}
	if value == 0 {
		return false, nil
	}
	if value == 1 {
		return true, nil
	}

	return false, fmt.Errorf("invalid value for boolean: %d (expected 0 or 1)", value)
}
