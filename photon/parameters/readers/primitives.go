package readers

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"golang.org/x/exp/constraints"
)

func readPrimitive[T constraints.Integer | constraints.Float](reader *bytes.Reader) (T, error) {
    var val T
    if err := binary.Read(reader, binary.BigEndian, &val); err != nil {
        return val, fmt.Errorf("failed to read %T: %w", val, err)
    }
    return val, nil
}


func ReadInt8(reader *bytes.Reader) (int8, error) {
    return readPrimitive[int8](reader)
}

func ReadInt16(reader *bytes.Reader) (int16, error) {
    return readPrimitive[int16](reader)
}

func ReadInt32(reader *bytes.Reader) (int32, error) {
    return readPrimitive[int32](reader)
}

func ReadInt64(reader *bytes.Reader) (int64, error) {
    return readPrimitive[int64](reader)
}

func ReadFloat32(reader *bytes.Reader) (float32, error) {
    return readPrimitive[float32](reader)
}

func ReadFloat64(reader *bytes.Reader) (float64, error) {
    return readPrimitive[float64](reader)
}

func ReadString(reader *bytes.Reader) (string, error){
	size, err := ReadInt16(reader)
	if err != nil{
		return "", err
	}

	buff := make([]byte, size)

	reader.Read(buff)

	return string(buff), nil
}

func ReadBoolean(readers *bytes.Reader) (bool, error){
	value, err := readPrimitive[uint8](readers)
	
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
