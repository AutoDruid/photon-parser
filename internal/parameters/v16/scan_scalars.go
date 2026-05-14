package v16

import (
	"math"

	"github.com/AutoDruid/photon-parser/internal/reader"
)

func scanString(reader *reader.Reader, value *Value) error {
	size, err := reader.ReadUInt16BE()
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

func scanFloat32(reader *reader.Reader, value *Value) error {
	res, err := reader.ReadFloat32BE()
	if err != nil {
		return err
	}
	value.Num = uint64(math.Float32bits(res))
	return nil
}

func scanFloat64(reader *reader.Reader, value *Value) error {
	res, err := reader.ReadFloat64BE()
	if err != nil {
		return err
	}
	value.Num = uint64(math.Float64bits(res))
	return nil
}

func scanInt8(reader *reader.Reader, value *Value) error {
	b, err := reader.ReadInt8()
	if err != nil {
		return err
	}
	value.Num = uint64(uint8(b))
	return nil
}

func scanInt16(reader *reader.Reader, value *Value) error {
	res, err := reader.ReadInt16BE()
	if err != nil {
		return err
	}
	value.Num = uint64(uint16(res))
	return nil
}

func scanInt32(reader *reader.Reader, value *Value) error {
	res, err := reader.ReadInt32BE()
	if err != nil {
		return err
	}
	value.Num = uint64(uint32(res))
	return nil
}

func scanInt64(reader *reader.Reader, value *Value) error {
	res, err := reader.ReadInt64BE()
	if err != nil {
		return err
	}
	value.Num = uint64(res)
	return nil
}

func scanBoolean(reader *reader.Reader, value *Value) error {
	bit, err := reader.ReadUInt8()
	if err != nil {
		return err
	}
	value.Num = uint64(bit)
	return nil
}
