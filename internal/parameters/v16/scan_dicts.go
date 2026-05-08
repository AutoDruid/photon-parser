package v16

import (
	"encoding/binary"
	"michelprogram/photon-parser/internal/reader"
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

	value.Blob, err = reader.ReadBytes(int(size))
	if err != nil {
		return err
	}

	value.Num = uint64(size)
	value.KeyType = ParameterType(keyType)
	value.ValType = ParameterType(valueType)
	
	return nil
}
