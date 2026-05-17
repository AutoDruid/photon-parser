package reader

import (
	"encoding/binary"
	"math"

	"github.com/AutoDruid/photon-parser/internal/errors"
)

// ReadInt8 reads an 8-bit signed integer from the reader.
// Returns an error if fewer than 1 byte is available.
func (r *Reader) ReadInt8() (int8, error) {
	size := r.Cursor + INT8_SIZE

	if size > r.Max {
		return 0, errors.ErrNotEnoughBytesInt8
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT8_SIZE
	return int8(b[0]), nil
}

// ReadUInt8 reads an 8-bit unsigned integer from the reader.
// Returns an error if fewer than 1 byte is available.
func (r *Reader) ReadUInt8() (uint8, error) {
	size := r.Cursor + INT8_SIZE

	if size > r.Max {
		return 0, errors.ErrNotEnoughBytesUInt8
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT8_SIZE
	return uint8(b[0]), nil
}

// ReadInt16BE reads a 16-bit signed integer from the reader in big endian byte order.
// Returns an error if fewer than 2 bytes are available.
func (r *Reader) ReadInt16BE() (int16, error) {

	size := r.Cursor + INT16_SIZE

	if size > r.Max {
		return 0, errors.ErrNotEnoughBytesInt16
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT16_SIZE

	return int16(binary.BigEndian.Uint16(b)), nil
}

// ReadInt16LE reads a 16-bit signed integer from the reader in little endian byte order.
// Returns an error if fewer than 2 bytes are available.
func (r *Reader) ReadInt16LE() (int16, error) {

	size := r.Cursor + INT16_SIZE

	if size > r.Max {
		return 0, errors.ErrNotEnoughBytesInt16
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT16_SIZE

	return int16(binary.LittleEndian.Uint16(b)), nil
}

// ReadUInt16BE reads a 16-bit unsigned integer from the reader in big endian byte order.
// Returns an error if fewer than 2 bytes are available.
func (r *Reader) ReadUInt16BE() (uint16, error) {

	size := r.Cursor + INT16_SIZE

	if size > r.Max {
		return 0, errors.ErrNotEnoughBytesUInt16
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT16_SIZE

	return binary.BigEndian.Uint16(b), nil
}

// ReadUInt16LE reads a 16-bit unsigned integer from the reader in little endian byte order.
// Returns an error if fewer than 2 bytes are available.
func (r *Reader) ReadUInt16LE() (uint16, error) {

	size := r.Cursor + INT16_SIZE

	if size > r.Max {
		return 0, errors.ErrNotEnoughBytesUInt16
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT16_SIZE

	return binary.LittleEndian.Uint16(b), nil
}

// ReadInt32BE reads a 32-bit signed integer from the reader in big endian byte order.
// Returns an error if fewer than 4 bytes are available.
func (r *Reader) ReadInt32BE() (int32, error) {
	size := r.Cursor + INT32_SIZE

	if size > r.Max {
		return 0, errors.ErrNotEnoughBytesInt32
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT32_SIZE

	return int32(binary.BigEndian.Uint32(b)), nil
}

// ReadInt32LE reads a 32-bit signed integer from the reader in little endian byte order.
// Returns an error if fewer than 4 bytes are available.
func (r *Reader) ReadInt32LE() (int32, error) {
	size := r.Cursor + INT32_SIZE

	if size > r.Max {
		return 0, errors.ErrNotEnoughBytesInt32
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT32_SIZE

	return int32(binary.LittleEndian.Uint32(b)), nil
}

// ReadUInt32BE reads a 32-bit unsigned integer from the reader in big endian byte order.
// Returns an error if fewer than 4 bytes are available.
func (r *Reader) ReadUInt32BE() (uint32, error) {
	size := r.Cursor + INT32_SIZE

	if size > r.Max {
		return 0, errors.ErrNotEnoughBytesUInt32
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT32_SIZE

	return binary.BigEndian.Uint32(b), nil
}

// ReadUInt32LE reads a 32-bit unsigned integer from the reader in little endian byte order.
// Returns an error if fewer than 4 bytes are available.
func (r *Reader) ReadUInt32LE() (uint32, error) {
	size := r.Cursor + INT32_SIZE

	if size > r.Max {
		return 0, errors.ErrNotEnoughBytesUInt32
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT32_SIZE

	return binary.LittleEndian.Uint32(b), nil
}

// ReadInt64BE reads a 64-bit signed integer from the reader in big endian byte order.
// Returns an error if fewer than 8 bytes are available.
func (r *Reader) ReadInt64BE() (int64, error) {
	size := r.Cursor + INT64_SIZE

	if size > r.Max {
		return 0, errors.ErrNotEnoughBytesInt64
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT64_SIZE

	return int64(binary.BigEndian.Uint64(b)), nil
}

// ReadInt64LE reads a 64-bit signed integer from the reader in little endian byte order.
// Returns an error if fewer than 8 bytes are available.
func (r *Reader) ReadInt64LE() (int64, error) {
	size := r.Cursor + INT64_SIZE

	if size > r.Max {
		return 0, errors.ErrNotEnoughBytesInt64
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT64_SIZE

	return int64(binary.LittleEndian.Uint64(b)), nil
}

// ReadUInt64BE reads a 64-bit unsigned integer from the reader in big endian byte order.
// Returns an error if fewer than 8 bytes are available.
func (r *Reader) ReadUInt64BE() (uint64, error) {
	size := r.Cursor + INT64_SIZE

	if size > r.Max {
		return 0, errors.ErrNotEnoughBytesUInt64
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT64_SIZE
	return binary.BigEndian.Uint64(b), nil
}

// ReadUInt64LE reads a 64-bit unsigned integer from the reader in little endian byte order.
// Returns an error if fewer than 8 bytes are available.
func (r *Reader) ReadUInt64LE() (uint64, error) {
	size := r.Cursor + INT64_SIZE

	if size > r.Max {
		return 0, errors.ErrNotEnoughBytesUInt64
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT64_SIZE
	return binary.LittleEndian.Uint64(b), nil
}

// ReadFloat32BE reads a 32-bit floating point number from the reader in big endian byte order.
// Returns an error if fewer than 4 bytes are available.
func (r *Reader) ReadFloat32BE() (float32, error) {
	size := r.Cursor + FLOAT32_SIZE

	if size > r.Max {
		return 0, errors.ErrNotEnoughBytesFloat32
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += FLOAT32_SIZE

	return math.Float32frombits(binary.BigEndian.Uint32(b)), nil
}

// ReadFloat32LE reads a 32-bit floating point number from the reader in little endian byte order.
// Returns an error if fewer than 4 bytes are available.
func (r *Reader) ReadFloat32LE() (float32, error) {
	size := r.Cursor + FLOAT32_SIZE

	if size > r.Max {
		return 0, errors.ErrNotEnoughBytesFloat32
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += FLOAT32_SIZE

	return math.Float32frombits(binary.LittleEndian.Uint32(b)), nil
}

// ReadFloat64BE reads a 64-bit floating point number from the reader in big endian byte order.
// Returns an error if fewer than 8 bytes are available.
func (r *Reader) ReadFloat64BE() (float64, error) {
	size := r.Cursor + FLOAT64_SIZE

	if size > r.Max {
		return 0, errors.ErrNotEnoughBytesFloat64
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT64_SIZE

	return math.Float64frombits(binary.BigEndian.Uint64(b)), nil
}

// ReadFloat64LE reads a 64-bit floating point number from the reader in little endian byte order.
// Returns an error if fewer than 8 bytes are available.
func (r *Reader) ReadFloat64LE() (float64, error) {
	size := r.Cursor + FLOAT64_SIZE

	if size > r.Max {
		return 0, errors.ErrNotEnoughBytesFloat64
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT64_SIZE

	return math.Float64frombits(binary.LittleEndian.Uint64(b)), nil
}

// ReadBoolean reads a boolean value from the reader.
// Format: single byte (0x00 = false, 0x01 = true).
// Returns an error if the value is not 0x00 or 0x01, or if no byte is available.
func (r *Reader) ReadBoolean() (bool, error) {
	bit, err := r.ReadUInt8()

	if err != nil {
		return false, err
	}

	if bit == 0 {
		return false, nil
	}

	if bit == 1 {
		return true, nil
	}

	return false, errors.ErrInvalidBooleanValue
}

// ReadString reads a string of n size from the reader.
// Returns an empty string if length is 0.
// Returns an error if the declared length cannot be fully read.
func (r *Reader) ReadString(n int) (string, error) {
	size := r.Cursor + n

	if size > r.Max {
		return "", errors.ErrNotEnoughBytesString
	}

	str := string(r.Buffer[r.Cursor:size])
	r.Cursor = size

	return str, nil
}

// ReadByte reads a single byte from the reader.
// Returns an error if fewer than 1 byte is available.
func (r *Reader) ReadByte() (byte, error) {

	size := r.Cursor + 1

	if size > r.Max {
		return 0, errors.ErrNotEnoughBytesByte
	}

	c := r.Buffer[r.Cursor]
	r.Cursor++
	return c, nil
}

// ReadBytes reads a []bytes of n size from the reader.
// Returns an empty []bytes if length is 0.
// Returns an error if the declared length cannot be fully read.
func (r *Reader) ReadBytes(n int) ([]byte, error) {

	size := r.Cursor + n

	if size > r.Max {
		return []byte{}, errors.ErrNotEnoughBytesBytes
	}

	buff := r.Buffer[r.Cursor:size]
	r.Cursor += n
	return buff, nil
}
