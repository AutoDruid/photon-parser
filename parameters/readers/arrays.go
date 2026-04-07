package readers

import (
	. "michelprogram/photon-parser/parser"

	"golang.org/x/exp/constraints"
)

func readArray[T constraints.Integer | constraints.Float](reader *Reader) ([]T, error) {
	size, err := ReadPrimitive[uint32](reader)

	if err != nil{
		return nil, err
	}

	val := make([]T, size)
	var i uint32 = 0

	for i = 0; i < size; i++{
		input, err := ReadPrimitive[T](reader)
  		if err != nil {
   			return nil, err
    	}
     	val[i] = input
	}
    return val, nil
}

func ReadInt8Array(reader *Reader) ([]int8, error){
	return readArray[int8](reader)
}

func ReadInt32Array(reader *Reader) ([]int32, error){
	return readArray[int32](reader)
}

func ReadStringArray(reader *Reader) ([]string, error){
	size, err := ReadPrimitive[uint32](reader)

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

func ReadArray(reader *Reader)([]any, error){

	size, err := ReadPrimitive[uint16](reader)

	if err != nil{
		return nil, err
	}

	ttype, err := ReadPrimitive[Type](reader)

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
