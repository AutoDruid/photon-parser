package errors

import "errors"

var HeaderSize = errors.New("header size too low")
var EncryptedPacket = errors.New("packet is encrypted or unknown: unexpected signature byte")  
var NotEnoughBytesInt8 = errors.New("not enough bytes to read int8")
var NotEnoughBytesUInt8 = errors.New("not enough bytes to read uint8")
var NotEnoughBytesInt16 = errors.New("not enough bytes to read int16")
var NotEnoughBytesUInt16 = errors.New("not enough bytes to read uint16")
var NotEnoughBytesInt32 = errors.New("not enough bytes to read int32")
var NotEnoughBytesUInt32 = errors.New("not enough bytes to read uint32")
var NotEnoughBytesInt64 = errors.New("not enough bytes to read int64")
var NotEnoughBytesUInt64 = errors.New("not enough bytes to read uint64")
var NotEnoughBytesFloat32 = errors.New("not enough bytes to read float32")
var NotEnoughBytesFloat64 = errors.New("not enough bytes to read float64")
var InvalidBooleanValue = errors.New("invalid value for boolean: (expected 0 or 1)")
var NotEnoughBytesString = errors.New("not enough bytes to read string")
var NotEnoughBytesByte = errors.New("not enough bytes to read byte")
var NotEnoughBytesBytes = errors.New("not enough bytes to read []byte")
var InvalidNegativeSkip = errors.New("negative skip")