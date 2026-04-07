package readers

import (
	"fmt"
	. "michelprogram/photon-parser/parser"
)

func ReadInt8(reader *Reader) (int8, error) {
    return ReadPrimitive[int8](reader)
}

func ReadInt16(reader *Reader) (int16, error) {
    return ReadPrimitive[int16](reader)
}

func ReadInt32(reader *Reader) (int32, error) {
    return ReadPrimitive[int32](reader)
}

func ReadInt64(reader *Reader) (int64, error) {
    return ReadPrimitive[int64](reader)
}

func ReadFloat32(reader *Reader) (float32, error) {
    return ReadPrimitive[float32](reader)
}

func ReadFloat64(reader *Reader) (float64, error) {
    return ReadPrimitive[float64](reader)
}

func ReadString(reader *Reader) (string, error){
	size, err := ReadInt16(reader)
	if err != nil{
		return "", err
	}

	buff := make([]byte, size)

	reader.Read(buff)

	return string(buff), nil
}

func ReadBoolean(readers *Reader) (bool, error){
	value, err := ReadPrimitive[uint8](readers)

	if err != nil{
		return false, err
	}
	if value == 0 {
		return false, nil
	}
	if value == 1 {
		return true, nil
	}

	return false, fmt.Errorf("Invalid value for boolean of %d", value)
}
