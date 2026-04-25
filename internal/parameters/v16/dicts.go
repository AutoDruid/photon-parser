package v16

import (
	"encoding/binary"
	"michelprogram/photon-parser/internal/reader"
)

// ReadDictionary reads a Photon Protocol16 dictionary with uniform key and value types.
// Format: Type byte (key type), Type byte (value type), uint16 size,
// then size key-value pairs where all keys have keyType and all values have valueType.
//
// Example wire format for map[int32]string{100: "hello"}:
//
//	0x69                 // keyType = Int32Type
//	0x73                 // valueType = StringType
//	0x00 0x01            // size = 1
//	0x00 0x00 0x00 0x64  // key: 100
//	0x00 0x05 'h' 'e' 'l' 'l' 'o'  // value: "hello"
//
// Returns an error if the dictionary cannot be fully read.
func (p Parameter) readDictionary(r *reader.Reader) (map[any]any, error) {

	keyType, err := r.ReadUInt8()
	if err != nil {
		return nil, err
	}

	valueType, err := r.ReadUInt8()
	if err != nil {
		return nil, err
	}

	size, err := r.ReadUInt16(binary.BigEndian)
	if err != nil {
		return nil, err
	}

	res := make(map[any]any, size)

	for i := uint16(0); i < size; i++ {
		key, err := p.decode(r, ParameterType(keyType))
		if err != nil {
			return nil, err
		}

		value, err := p.decode(r, ParameterType(valueType))
		if err != nil {
			return nil, err
		}

		res[key] = value
	}

	return res, nil
}

// ReadHashTable reads a Photon Protocol16 hashtable with mixed types.
// Format: uint16 size, then size key-value pairs where each pair has its own
// type bytes (keyType, key, valueType, value).
//
// This differs from Dictionary in that each entry can have different types,
// whereas Dictionary requires all keys to be the same type and all values
// to be the same type.
//
// Example wire format for map{100: "hello", "key": true}:
//
//	0x00 0x02            // size = 2
//	0x69                 // entry 1 key type = Int32Type
//	0x00 0x00 0x00 0x64  // entry 1 key: 100
//	0x73                 // entry 1 value type = StringType
//	0x00 0x05 'h' 'e' 'l' 'l' 'o'  // entry 1 value: "hello"
//	0x73                 // entry 2 key type = StringType
//	0x00 0x03 'k' 'e' 'y'  // entry 2 key: "key"
//	0x6f                 // entry 2 value type = BooleanType
//	0x01                 // entry 2 value: true
//
// Returns an error if the hashtable cannot be fully read.
func (p Parameter) readHashTable(r *reader.Reader) (map[any]any, error) {
	size, err := r.ReadUInt16(binary.BigEndian)
	if err != nil {
		return nil, err
	}

	res := make(map[any]any, int(size))
	for i := uint16(0); i < size; i++ {
		keyType, err := r.ReadUInt8()
		if err != nil {
			return nil, err
		}
		key, err := p.decode(r, ParameterType(keyType))
		if err != nil {
			return nil, err
		}

		valueType, err := r.ReadUInt8()
		if err != nil {
			return nil, err
		}
		value, err := p.decode(r, ParameterType(valueType))
		if err != nil {
			return nil, err
		}

		res[key] = value
	}

	return res, nil
}
