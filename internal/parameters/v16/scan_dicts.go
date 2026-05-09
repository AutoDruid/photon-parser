package v16

import (
	"encoding/binary"

	"github.com/AutoDruid/photon-parser/internal/reader"
)

func scanDictionary(reader *reader.Reader, value *Value) error {
	keyType, err := reader.ReadUInt8()
	if err != nil {
		return err
	}

	valueType, err := reader.ReadUInt8()
	if err != nil {
		return err
	}

	size, err := reader.ReadUInt16(binary.BigEndian)
	if err != nil {
		return err
	}

	start := reader.Cursor
	for i := uint16(0); i < size; i++ {
		if _, err := scanPayload(reader, ParameterType(keyType)); err != nil {
			return err
		}
		if _, err := scanPayload(reader, ParameterType(valueType)); err != nil {
			return err
		}
	}
	value.Blob = reader.Buffer[start:reader.Cursor]
	value.Num = uint64(size)
	value.KeyType = ParameterType(keyType)
	value.ValType = ParameterType(valueType)

	return nil
}
