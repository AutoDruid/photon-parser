package reader

import (
	"encoding/binary"
	"fmt"
	"math"
	"michelprogram/photon-parser/internal/hooks"
	"michelprogram/photon-parser/internal/types"
)

// ParameterParser is implemented by each protocol-version parameters package
// (v16, v18, ...). It is wired once at Parser construction so the hot path
// has no version branches.
type ParameterParser interface {
	Parse(r *Reader, out *types.Parameter, hooks *hooks.Hooks) error
}

type ReliableHeaderParameterCount interface {
	Count(r *Reader) (int, error)
}

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

type Options struct {
	ParameterParser
	ReliableHeaderParameterCount
	BinaryOrder binary.ByteOrder
}

type Reader struct {
	Buffer []byte
	Max    int
	Cursor int

	Options
}

func NewReader(data []byte, options Options) *Reader {
	return &Reader{
		Buffer:  data,
		Max:     len(data),
		Cursor:  0,
		Options: options,
	}
}

func (r *Reader) SetParameterParser(parser ParameterParser) {
	r.ParameterParser = parser
}

// ReadInt8 reads an 8-bit signed integer from the reader.
// Returns an error if fewer than 1 byte is available.
func (r *Reader) ReadInt8() (int8, error) {
	size := r.Cursor + INT8_SIZE

	if size > r.Max {
		return 0, fmt.Errorf("not enough bytes to read int8")
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
		return 0, fmt.Errorf("not enough bytes to read uint8")
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT8_SIZE
	return uint8(b[0]), nil
}

// ReadInt16 reads a 16-bit signed integer from the reader in big-endian format.
// Returns an error if fewer than 2 bytes are available.
func (r *Reader) ReadInt16() (int16, error) {

	size := r.Cursor + INT16_SIZE

	if size > r.Max {
		return 0, fmt.Errorf("not enough bytes to read int16")
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT16_SIZE

	return int16(b[0])<<8 | int16(b[1]), nil
}

func (r *Reader) ReadInt16LittleEndian() (int16, error) {

	size := r.Cursor + INT16_SIZE

	if size > r.Max {
		return 0, fmt.Errorf("not enough bytes to read int16")
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT16_SIZE

	return int16(b[1])<<8 | int16(b[0]), nil
}

// ReadUInt16 reads a 16-bit unsigned integer from the reader in big-endian format.
// Returns an error if fewer than 2 bytes are available.
func (r *Reader) ReadUInt16() (uint16, error) {

	size := r.Cursor + INT16_SIZE

	if size > r.Max {
		return 0, fmt.Errorf("not enough bytes to read uint16")
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT16_SIZE

	return uint16(b[0])<<8 | uint16(b[1]), nil
}

func (r *Reader) ReadUInt16LittleEndian() (uint16, error) {

	size := r.Cursor + INT16_SIZE

	if size > r.Max {
		return 0, fmt.Errorf("not enough bytes to read uint16")
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT16_SIZE

	return uint16(b[1])<<8 | uint16(b[0]), nil
}

// ReadInt32 reads a 32-bit signed integer from the reader in big-endian format.
// Returns an error if fewer than 4 bytes are available.
func (r *Reader) ReadInt32() (int32, error) {
	size := r.Cursor + INT32_SIZE

	if size > r.Max {
		return 0, fmt.Errorf("not enough bytes to read int32")
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT32_SIZE

	return int32(b[0])<<24 |
		int32(b[1])<<16 |
		int32(b[2])<<8 |
		int32(b[3]), nil
}

func (r *Reader) ReadInt32LittleEndian() (int32, error) {
	size := r.Cursor + INT32_SIZE

	if size > r.Max {
		return 0, fmt.Errorf("not enough bytes to read int32")
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT32_SIZE

	return int32(b[3])<<24 |
		int32(b[2])<<16 |
		int32(b[1])<<8 |
		int32(b[0]), nil
}

// ReadUInt32 reads a 32-bit unsigned integer from the reader in big-endian format.
// Returns an error if fewer than 4 bytes are available.
func (r *Reader) ReadUInt32() (uint32, error) {
	size := r.Cursor + INT32_SIZE

	if size > r.Max {
		return 0, fmt.Errorf("not enough bytes to read uint32")
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT32_SIZE

	return uint32(b[0])<<24 |
		uint32(b[1])<<16 |
		uint32(b[2])<<8 |
		uint32(b[3]), nil
}

// ReadInt64 reads a 64-bit signed integer from the reader in big-endian format.
// Returns an error if fewer than 8 bytes are available.
func (r *Reader) ReadInt64() (int64, error) {
	size := r.Cursor + INT64_SIZE

	if size > r.Max {
		return 0, fmt.Errorf("not enough bytes to read int64")
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT64_SIZE

	return int64(b[0])<<56 |
		int64(b[1])<<48 |
		int64(b[2])<<40 |
		int64(b[3])<<32 |
		int64(b[4])<<24 |
		int64(b[5])<<16 |
		int64(b[6])<<8 |
		int64(b[7]), nil
}

// ReadUInt64 reads a 64-bit unsigned integer from the reader in big-endian format.
// Returns an error if fewer than 8 bytes are available.
func (r *Reader) ReadUInt64() (uint64, error) {
	size := r.Cursor + INT64_SIZE

	if size > r.Max {
		return 0, fmt.Errorf("not enough bytes to read uint64")
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT64_SIZE
	return uint64(b[0])<<56 |
		uint64(b[1])<<48 |
		uint64(b[2])<<40 |
		uint64(b[3])<<32 |
		uint64(b[4])<<24 |
		uint64(b[5])<<16 |
		uint64(b[6])<<8 |
		uint64(b[7]), nil
}

// ReadFloat32 reads a 32-bit floating point number from the reader in big-endian format.
// Returns an error if fewer than 4 bytes are available.
func (r *Reader) ReadFloat32() (float32, error) {
	size := r.Cursor + FLOAT32_SIZE

	if size > r.Max {
		return 0, fmt.Errorf("not enough bytes to read float32")
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += FLOAT32_SIZE

	return math.Float32frombits(r.BinaryOrder.Uint32(b)), nil
}

// ReadFloat64 reads a 64-bit floating point number from the reader in big-endian format.
// Returns an error if fewer than 8 bytes are available.
func (r *Reader) ReadFloat64() (float64, error) {
	size := r.Cursor + FLOAT64_SIZE

	if size > r.Max {
		return 0, fmt.Errorf("not enough bytes to read float64")
	}

	b := r.Buffer[r.Cursor:size]
	r.Cursor += INT64_SIZE

	return math.Float64frombits(r.BinaryOrder.Uint64(b)), nil
}

// ReadBoolean reads a boolean value from the reader.
// Format: single byte (0x00 = false, 0x01 = true).
// Returns an error if the value is not 0x00 or 0x01, or if no byte is available.
func (r *Reader) ReadBoolean() (bool, error) {
	bit, err := r.ReadUInt8()

	if err != nil {
		return false, fmt.Errorf("not enough bytes to read boolean")
	}

	if bit == 0 {
		return false, nil
	}

	if bit == 1 {
		return true, nil
	}

	return false, fmt.Errorf("invalid value for boolean: %d (expected 0 or 1)", bit)
}

// ReadString reads a string of n size from the reader.
// Returns an empty string if length is 0.
// Returns an error if the declared length cannot be fully read.
func (r *Reader) ReadString(n int) (string, error) {
	size := r.Cursor + n

	if size > r.Max {
		return "", fmt.Errorf("not enough bytes to read string of size %d", n)
	}

	str := string(r.Buffer[r.Cursor:size])
	r.Cursor = size

	return str, nil
}

func (r *Reader) ReadByte() (byte, error) {

	size := r.Cursor + 1

	if size > r.Max {
		return 0, fmt.Errorf("not enough bytes to read byte")
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
		return []byte{}, fmt.Errorf("not enough bytes to read []byte of size %d", n)
	}

	buff := r.Buffer[r.Cursor:size]
	r.Cursor += n
	return buff, nil
}

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

func (r *Reader) ReadVarintInt32() (int32, error) {

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

	//ZigZag decode
	return int32((res >> 1) ^ uint32(-(int32(res & 1)))), nil
}

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

func (r *Reader) ReadVarintInt64() (int64, error) {

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

	//ZigZag decode
	return int64((res >> 1) ^ uint64(-(int64(res & 1)))), nil
}

func (r *Reader) Reset(data []byte) {
	r.Buffer = data
	r.Cursor = 0
	r.Max = len(data)
}
