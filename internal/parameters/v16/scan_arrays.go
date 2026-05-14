package v16

import (
	"github.com/AutoDruid/photon-parser/internal/reader"
)

func scanInt8Array(reader *reader.Reader, value *Value) error {
	size, err := reader.ReadUInt32BE()
	if err != nil {
		return err
	}

	value.Blob, err = reader.ReadBytes(int(size))
	if err != nil {
		return err
	}

	value.Num = uint64(size)
	return nil
}

func scanInt32Array(reader *reader.Reader, value *Value) error {
	size, err := reader.ReadUInt32BE()
	if err != nil {
		return err
	}

	value.Blob, err = reader.ReadBytes(int(size * 4))
	if err != nil {
		return err
	}
	value.Num = uint64(size)
	return nil
}

func scanStringArray(reader *reader.Reader, value *Value) error {
	size, err := reader.ReadUInt32BE()
	if err != nil {
		return err
	}

	start := reader.Cursor
	for i := uint32(0); i < size; i++ {
		size, err := reader.ReadUInt16BE()
		if err != nil {
			return err
		}

		err = reader.Skip(int(size))
		if err != nil {
			return err
		}
	}
	value.Blob = reader.Buffer[start:reader.Cursor]
	value.Num = uint64(size)
	return nil
}

func scanArray(reader *reader.Reader, value *Value) error {
	size, err := reader.ReadUInt16BE()
	if err != nil {
		return err
	}

	ttype, err := reader.ReadUInt8()
	if err != nil {
		return err
	}

	value.KeyType = ParameterType(ttype)

	start := reader.Cursor
	for i := uint16(0); i < size; i++ {
		_, err = scanPayload(reader, value.KeyType)
		if err != nil {
			return err
		}
	}

	value.Blob = reader.Buffer[start:reader.Cursor]
	if err != nil {
		return err
	}
	value.Num = uint64(size)
	return nil
}
