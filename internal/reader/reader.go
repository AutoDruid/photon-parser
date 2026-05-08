package reader

import (
	"encoding/binary"
	"math"
	"michelprogram/photon-parser/internal/errors"
)

const (
	INT8_SIZE       = 1
	INT16_SIZE      = 2
	INT32_SIZE      = 4
	INT64_SIZE      = 8
	FLOAT32_SIZE    = 4
	FLOAT64_SIZE    = 8
	VARINT_SHIFT    = 7
	VARINT_MASK     = 0x7F
	VARINT_MSB_MASK = 0x80
)

type Reader struct {
	Buffer []byte
	Max    int
	Cursor int
}

func NewReader(data []byte) *Reader {
	return &Reader{
		Buffer: data,
		Max:    len(data),
		Cursor: 0,
	}
}

// Read Remaining bytes from the reader.
func (r *Reader) ReadRemaining() []byte {
	tmp := r.Cursor
	r.Cursor = r.Max
	return r.Buffer[tmp:]
}

func (r *Reader) Skip(n int) error {

	if n < 0 {
		return errors.InvalidNegativeSkip
	}

	size := r.Cursor + n

	if size > r.Max {
		return errors.NotEnoughBytesBytes
	}

	r.Cursor += n
	return nil
}
// ReadInt8 reads an 8-bit signed integer from the reader.
// Returns an error if fewer than 1 byte is available.
func (r *Reader) ReadInt8() (int8, error) {
	size := r.Cursor + INT8_SIZE

	if size > r.Max {
		return 0, errors.NotEnoughBytesInt8
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
		return 0, errors.NotEnoughBytesUInt8
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT8_SIZE
	return uint8(b[0]), nil
}

// ReadInt16 reads a 16-bit signed integer from the reader in the given byte order.
// Returns an error if fewer than 2 bytes are available.
func (r *Reader) ReadInt16(order binary.ByteOrder) (int16, error) {

	size := r.Cursor + INT16_SIZE

	if size > r.Max {
		return 0, errors.NotEnoughBytesInt16
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT16_SIZE

	return int16(order.Uint16(b)), nil
}

// ReadUInt16 reads a 16-bit unsigned integer from the reader in the given byte order.
// Returns an error if fewer than 2 bytes are available.
func (r *Reader) ReadUInt16(order binary.ByteOrder) (uint16, error) {

	size := r.Cursor + INT16_SIZE

	if size > r.Max {
		return 0, errors.NotEnoughBytesUInt16
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT16_SIZE

	return order.Uint16(b), nil
}

// ReadInt32 reads a 32-bit signed integer from the reader in the given byte order.
// Returns an error if fewer than 4 bytes are available.
func (r *Reader) ReadInt32(order binary.ByteOrder) (int32, error) {
	size := r.Cursor + INT32_SIZE

	if size > r.Max {
		return 0, errors.NotEnoughBytesInt32
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT32_SIZE

	return int32(order.Uint32(b)), nil
}

// ReadUInt32 reads a 32-bit unsigned integer from the reader in the given byte order.
// Returns an error if fewer than 4 bytes are available.
func (r *Reader) ReadUInt32(order binary.ByteOrder) (uint32, error) {
	size := r.Cursor + INT32_SIZE

	if size > r.Max {
		return 0, errors.NotEnoughBytesUInt32
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT32_SIZE

	return order.Uint32(b), nil
}

// ReadInt64 reads a 64-bit signed integer from the reader in the given byte order.
// Returns an error if fewer than 8 bytes are available.
func (r *Reader) ReadInt64(order binary.ByteOrder) (int64, error) {
	size := r.Cursor + INT64_SIZE

	if size > r.Max {
		return 0, errors.NotEnoughBytesInt64
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT64_SIZE

	return int64(order.Uint64(b)), nil
}

// ReadUInt64 reads a 64-bit unsigned integer from the reader in the given byte order.
// Returns an error if fewer than 8 bytes are available.
func (r *Reader) ReadUInt64(order binary.ByteOrder) (uint64, error) {
	size := r.Cursor + INT64_SIZE

	if size > r.Max {
		return 0, errors.NotEnoughBytesUInt64
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT64_SIZE
	return order.Uint64(b), nil
}

// ReadFloat32 reads a 32-bit floating point number from the reader in the given byte order.
// Returns an error if fewer than 4 bytes are available.
func (r *Reader) ReadFloat32(order binary.ByteOrder) (float32, error) {
	size := r.Cursor + FLOAT32_SIZE

	if size > r.Max {
		return 0, errors.NotEnoughBytesFloat32
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += FLOAT32_SIZE

	return math.Float32frombits(order.Uint32(b)), nil
}

// ReadFloat64 reads a 64-bit floating point number from the reader in the given byte order.
// Returns an error if fewer than 8 bytes are available.
func (r *Reader) ReadFloat64(order binary.ByteOrder) (float64, error) {
	size := r.Cursor + FLOAT64_SIZE

	if size > r.Max {
		return 0, errors.NotEnoughBytesFloat64
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT64_SIZE

	return math.Float64frombits(order.Uint64(b)), nil
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

	return false, errors.InvalidBooleanValue
}

// ReadString reads a string of n size from the reader.
// Returns an empty string if length is 0.
// Returns an error if the declared length cannot be fully read.
func (r *Reader) ReadString(n int) (string, error) {
	size := r.Cursor + n

	if size > r.Max {
		return "", errors.NotEnoughBytesString
	}

	str := string(r.Buffer[r.Cursor:size])
	//str := unsafe.String(&r.Buffer[r.Cursor], n)
	r.Cursor = size

	return str, nil
}

func (r *Reader) ReadByte() (byte, error) {

	size := r.Cursor + 1

	if size > r.Max {
		return 0, errors.NotEnoughBytesByte
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
		return []byte{}, errors.NotEnoughBytesBytes
	}

	buff := r.Buffer[r.Cursor:size]
	r.Cursor += n
	return buff, nil
}

// ReadVarintUInt32 reads a 32-bit unsigned integer from the reader in varint format.
func (r *Reader) ReadVarintUInt32() (uint32, error) {

	var res uint32
	var shift uint8
	for {

		buff, err := r.ReadByte()
		if err != nil {
			return 0, err
		}

		res |= uint32(buff&VARINT_MASK) << shift

		if buff&VARINT_MSB_MASK == 0 {
			break
		}
		shift += VARINT_SHIFT
	}

	return res, nil
}

// ReadVarintInt32 reads a 32-bit signed integer from the reader in varint format.
func (r *Reader) ReadVarintInt32() (int32, error) {

	res, err := r.ReadVarintUInt32()
	if err != nil {
		return 0, err
	}

	//ZigZag decode
	return int32((res >> 1) ^ uint32(-(int32(res & 1)))), nil
}

// ReadVarintUInt64 reads a 64-bit unsigned integer from the reader in varint format.
func (r *Reader) ReadVarintUInt64() (uint64, error) {

	var res uint64
	var shift uint8
	for {

		buff, err := r.ReadByte()
		if err != nil {
			return 0, err
		}

		res |= uint64(buff&VARINT_MASK) << shift

		if buff&VARINT_MSB_MASK == 0 {
			break
		}
		shift += VARINT_SHIFT
	}

	return res, nil
}

// ReadVarintInt64 reads a 64-bit signed integer from the reader in varint format.
func (r *Reader) ReadVarintInt64() (int64, error) {

	res, err := r.ReadVarintUInt64()
	if err != nil {
		return 0, err
	}

	//ZigZag decode
	return int64((res >> 1) ^ uint64(-(int64(res & 1)))), nil
}

// Reset resets the reader to the beginning of the buffer.
func (r *Reader) Reset(data []byte) {
	r.Buffer = data
	r.Cursor = 0
	r.Max = len(data)
}
