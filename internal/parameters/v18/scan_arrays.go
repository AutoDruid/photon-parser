package v18

import (
	"github.com/AutoDruid/photon-parser/internal/reader"
)

func scanFloat32Array(reader *reader.Reader, value *Value) error {
	count, err := reader.ReadVarintUInt32()
	if err != nil {
		return err
	}
	blob, err := reader.ReadBytes(int(count) * 4)
	if err != nil {
		return err
	}
	value.Blob = blob
	value.Num = uint64(count)
	return nil
}

func scanShortArray(reader *reader.Reader, value *Value) error {
	count, err := reader.ReadVarintUInt32()
	if err != nil {
		return err
	}
	value.Blob, err = reader.ReadBytes(int(count) * 2)
	if err != nil {
		return err
	}
	value.Num = uint64(count)
	return nil
}

func scanByteArray(reader *reader.Reader, value *Value) error {
	count, err := reader.ReadVarintUInt32()
	if err != nil {
		return err
	}
	value.Blob, err = reader.ReadBytes(int(count))
	if err != nil {
		return err
	}
	value.Num = uint64(count)
	return nil
}

func scanBooleanArray(reader *reader.Reader, value *Value) error {
	count, err := reader.ReadVarintUInt32()
	if err != nil {
		return err
	}
	packedBytes := (int(count) + 7) / 8
	value.Blob, err = reader.ReadBytes(packedBytes)
	if err != nil {
		return err
	}
	value.Num = uint64(count)
	return nil
}

func scanStringArray(reader *reader.Reader, value *Value) error {
	count, err := reader.ReadVarintUInt32()
	if err != nil {
		return err
	}
	start := reader.Cursor
	for i := uint32(0); i < count; i++ {
		size, err := reader.ReadVarintUInt32()
		if err != nil {
			return err
		}
		if _, err := reader.ReadBytes(int(size)); err != nil {
			return err
		}
	}
	value.Blob = reader.Buffer[start:reader.Cursor]
	value.Num = uint64(count)
	return nil
}

func scanCompressedIntArray(reader *reader.Reader, value *Value) error {
	count, err := reader.ReadVarintUInt32()
	if err != nil {
		return err
	}
	start := reader.Cursor
	for i := uint32(0); i < count; i++ {
		if _, err := reader.ReadVarintInt32(); err != nil {
			return err
		}
	}
	value.Blob = reader.Buffer[start:reader.Cursor]
	value.Num = uint64(count)
	return nil
}

func scanCompressedLongArray(reader *reader.Reader, value *Value) error {
	count, err := reader.ReadVarintUInt32()
	if err != nil {
		return err
	}
	start := reader.Cursor
	for i := uint32(0); i < count; i++ {
		if _, err := reader.ReadVarintInt64(); err != nil {
			return err
		}
	}
	value.Blob = reader.Buffer[start:reader.Cursor]
	value.Num = uint64(count)
	return nil
}

func scanArray(r *reader.Reader, value *Value) error {
	count, err := r.ReadVarintUInt32()
	if err != nil {
		return err
	}
	start := r.Cursor
	for i := uint32(0); i < count; i++ {
		ttype, err := r.ReadByte()
		if err != nil {
			return err
		}
		if _, err := scanPayload(r, ParameterType(ttype)); err != nil {
			return err
		}
	}
	value.Blob = r.Buffer[start:r.Cursor]
	value.Num = uint64(count)
	return nil
}
