package v16

import (
	"encoding/binary"
	"michelprogram/photon-parser/internal/reader"
)

func scanInt8Array(reader *reader.Reader, value *Value) error {
	size, err := reader.ReadUInt32(binary.BigEndian)
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
	size, err := reader.ReadUInt32(binary.BigEndian)
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

func scanStringArray(reader *reader.Reader, value *Value) error {
	size, err := reader.ReadUInt32(binary.BigEndian)
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

func scanArray(reader *reader.Reader, value *Value) error {
	size, err := reader.ReadUInt32(binary.BigEndian)
	if err != nil {
		return err
	}

	value.Blob, err = reader.ReadBytes(int(size+1))
	if err != nil {
		return err
	}
	value.Num = uint64(size)
	return nil
}