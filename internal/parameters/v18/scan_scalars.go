package v18

import (
	"encoding/binary"
	"math"
	"michelprogram/photon-parser/internal/reader"
)

func scanString(reader *reader.Reader, value *Value) error {
	size, err := reader.ReadVarintUInt32()
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
	res, err := reader.ReadFloat32(binary.LittleEndian)
	if err != nil {
		return err
	}
	value.Num = uint64(math.Float32bits(res))
	return nil
}

func scanInt8(reader *reader.Reader, value *Value) error {
	b, err := reader.ReadByte()
	if err != nil {
		return err
	}
	value.Num = uint64(uint8(b))
	return nil
}

func scanInt8Positive(reader *reader.Reader, value *Value) error {
	b, err := reader.ReadByte()
	if err != nil {
		return err
	}
	value.Num = uint64(uint(int32(b)))
	return nil
}

func scanInt8Negative(reader *reader.Reader, value *Value) error {
	b, err := reader.ReadByte()
	if err != nil {
		return err
	}
	value.Num = uint64(uint32(-int32(b)))
	return nil
}

func scanInt16Type(reader *reader.Reader, value *Value) error {
	res, err := reader.ReadInt16(binary.LittleEndian)
	if err != nil {
		return err
	}
	value.Num = uint64(uint16(res))
	return nil
}

func scanInt16Positive(reader *reader.Reader, value *Value) error {
	res, err := reader.ReadUInt16(binary.LittleEndian)
	if err != nil {
		return err
	}
	value.Num = uint64(uint32(int32(res)))
	return nil
}

func scanInt16Negative(reader *reader.Reader, value *Value) error {
	res, err := reader.ReadUInt16(binary.LittleEndian)
	if err != nil {
		return err
	}
	value.Num = uint64(uint32(-int32(res)))
	return nil
}

func scanLong8Positive(reader *reader.Reader, value *Value) error {
	res, err := reader.ReadByte()
	if err != nil {
		return err
	}
	value.Num = uint64(res)
	return nil
}

func scanLong8Negative(reader *reader.Reader, value *Value) error {
	res, err := reader.ReadByte()
	if err != nil {
		return err
	}
	value.Num = uint64(-int64(res))
	return nil
}

func scanLong16Positive(reader *reader.Reader, value *Value) error {
	res, err := reader.ReadUInt16(binary.LittleEndian)
	if err != nil {
		return err
	}
	value.Num = uint64(res)
	return nil
}

func scanLong16Negative(reader *reader.Reader, value *Value) error {
	res, err := reader.ReadUInt16(binary.LittleEndian)
	if err != nil {
		return err
	}
	value.Num = uint64(-int64(res))
	return nil
}

func scanCompressedInt32(reader *reader.Reader, value *Value) error {
	res, err := reader.ReadVarintInt32()
	if err != nil {
		return err
	}
	value.Num = uint64(uint32(res))
	return nil
}

func scanCompressedInt64(reader *reader.Reader, value *Value) error {
	res, err := reader.ReadVarintInt64()
	if err != nil {
		return err
	}
	value.Num = uint64(res)
	return nil
}
