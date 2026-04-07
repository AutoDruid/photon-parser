package readers

import (
	. "michelprogram/photon-parser/parser"
)

func ReadDictionnary(reader *Reader)(map[any]any, error){

	keyType, err := ReadPrimitive[Type](reader)

	if err != nil{
		return nil, err
	}

	valueType, err := ReadPrimitive[Type](reader)

	if err != nil{
		return nil, err
	}

	size, err := ReadPrimitive[uint16](reader)

	if err != nil{
		return nil, err
	}

	res := make(map[any]any, size)
	var i uint16 = 0

	for i = 0; i < size; i++{
		key, err := Decode(reader, keyType)
  		if err != nil {
   			return nil, err
    	}

  		value, err := Decode(reader, valueType)
    	if err != nil {
      		return nil, err
       	}

        res[key] = value
	}

	return res, nil
}

func ReadHashtable(reader *Reader) (map[any]any, error) {
	size, err := ReadPrimitive[uint16](reader)
	if err != nil {
		return nil, err
	}


	res := make(map[any]any, int(size))
	for i := uint16(0); i < size; i++ {
		keyType, err := ReadPrimitive[Type](reader)
		if err != nil {
			return nil, err
		}
		key, err := Decode(reader, keyType)
		if err != nil {
			return nil, err
		}

		valueType, err := ReadPrimitive[Type](reader)
		if err != nil {
			return nil, err
		}
		value, err := Decode(reader, valueType)
		if err != nil {
			return nil, err
		}

		res[key] = value
	}

	return res, nil
}
