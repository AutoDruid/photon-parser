package errors

import "errors"

var ErrHeaderSize = errors.New("header size too low")
var ErrEncryptedPacket = errors.New("packet is encrypted or unknown: unexpected signature byte")
var ErrNotEnoughBytesInt8 = errors.New("not enough bytes to read int8")
var ErrNotEnoughBytesUInt8 = errors.New("not enough bytes to read uint8")
var ErrNotEnoughBytesInt16 = errors.New("not enough bytes to read int16")
var ErrNotEnoughBytesUInt16 = errors.New("not enough bytes to read uint16")
var ErrNotEnoughBytesInt32 = errors.New("not enough bytes to read int32")
var ErrNotEnoughBytesUInt32 = errors.New("not enough bytes to read uint32")
var ErrNotEnoughBytesInt64 = errors.New("not enough bytes to read int64")
var ErrNotEnoughBytesUInt64 = errors.New("not enough bytes to read uint64")
var ErrNotEnoughBytesFloat32 = errors.New("not enough bytes to read float32")
var ErrNotEnoughBytesFloat64 = errors.New("not enough bytes to read float64")
var ErrInvalidBooleanValue = errors.New("invalid value for boolean: (expected 0 or 1)")
var ErrNotEnoughBytesString = errors.New("not enough bytes to read string")
var ErrNotEnoughBytesByte = errors.New("not enough bytes to read byte")
var ErrNotEnoughBytesBytes = errors.New("not enough bytes to read []byte")
var ErrInvalidNegativeSkip = errors.New("negative skip")
