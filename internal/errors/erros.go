package errors

import "errors"

// ErrHeaderSize indicates a command header length smaller than the minimum header size.
var ErrHeaderSize = errors.New("header size too low")

// ErrEncryptedPacket indicates a reliable packet that cannot be decoded as plain Photon payload.
var ErrEncryptedPacket = errors.New("packet is encrypted or unknown: unexpected signature byte")

// ErrNotEnoughBytesInt8 indicates an underflow while reading an int8 value.
var ErrNotEnoughBytesInt8 = errors.New("not enough bytes to read int8")

// ErrNotEnoughBytesUInt8 indicates an underflow while reading a uint8 value.
var ErrNotEnoughBytesUInt8 = errors.New("not enough bytes to read uint8")

// ErrNotEnoughBytesInt16 indicates an underflow while reading an int16 value.
var ErrNotEnoughBytesInt16 = errors.New("not enough bytes to read int16")

// ErrNotEnoughBytesUInt16 indicates an underflow while reading a uint16 value.
var ErrNotEnoughBytesUInt16 = errors.New("not enough bytes to read uint16")

// ErrNotEnoughBytesInt32 indicates an underflow while reading an int32 value.
var ErrNotEnoughBytesInt32 = errors.New("not enough bytes to read int32")

// ErrNotEnoughBytesUInt32 indicates an underflow while reading a uint32 value.
var ErrNotEnoughBytesUInt32 = errors.New("not enough bytes to read uint32")

// ErrNotEnoughBytesInt64 indicates an underflow while reading an int64 value.
var ErrNotEnoughBytesInt64 = errors.New("not enough bytes to read int64")

// ErrNotEnoughBytesUInt64 indicates an underflow while reading a uint64 value.
var ErrNotEnoughBytesUInt64 = errors.New("not enough bytes to read uint64")

// ErrNotEnoughBytesFloat32 indicates an underflow while reading a float32 value.
var ErrNotEnoughBytesFloat32 = errors.New("not enough bytes to read float32")

// ErrNotEnoughBytesFloat64 indicates an underflow while reading a float64 value.
var ErrNotEnoughBytesFloat64 = errors.New("not enough bytes to read float64")

// ErrInvalidBooleanValue indicates a boolean byte outside the supported 0 or 1 values.
var ErrInvalidBooleanValue = errors.New("invalid value for boolean: (expected 0 or 1)")

// ErrNotEnoughBytesString indicates an underflow while reading a string value.
var ErrNotEnoughBytesString = errors.New("not enough bytes to read string")

// ErrNotEnoughBytesByte indicates an underflow while reading one byte.
var ErrNotEnoughBytesByte = errors.New("not enough bytes to read byte")

// ErrNotEnoughBytesBytes indicates an underflow while reading a byte slice.
var ErrNotEnoughBytesBytes = errors.New("not enough bytes to read []byte")

// ErrInvalidNegativeSkip indicates an attempt to skip a negative number of bytes.
var ErrInvalidNegativeSkip = errors.New("negative skip")
