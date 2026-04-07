package readers

import (
	"bytes"

	"golang.org/x/exp/constraints"
)

func readArray[T constraints.Integer | constraints.Float](reader *bytes.Reader) ([]T, error) {
	size, err := readPrimitive[uint32](reader)

	if err != nil{
		return nil, err
	}

	val := make([]T, size)
	var i uint32 = 0

	for i = 0; i < size; i++{
		input, err := readPrimitive[T](reader)
  		if err != nil {
   			return nil, err
    	}
     	val[i] = input
	}
    return val, nil
}

func ReadInt8Array(reader *bytes.Reader) ([]int8, error){
	return readArray[int8](reader)
}

func ReadInt32Array(reader *bytes.Reader) ([]int32, error){
	return readArray[int32](reader)
}

func ReadStringArray(reader *bytes.Reader) ([]string, error){
	size, err := readPrimitive[uint32](reader)

	if err != nil{
		return nil, err
	}

	val := make([]string, size)
	var i uint32 = 0

	for i = 0; i < size; i++{
		input, err := ReadString(reader)
  		if err != nil {
   			return nil, err
    	}
     	val[i] = input
	}
    return val, nil
}

func ReadArray(reader *bytes.Reader)([]any, error){

	size, err := readPrimitive[uint16](reader)

	if err != nil{
		return nil, err
	}

	ttype, err := readPrimitive[Type](reader)

	if err != nil{
		return nil, err
	}

	val := make([]any, size)
	var i uint16 = 0

	for i = 0; i < size; i++{
		input, err := Decode(reader, ttype)
		if err != nil{
			return nil, err
		}
     	val[i] = input
      }
    return val, nil
}
