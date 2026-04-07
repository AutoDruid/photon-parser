package parameters

import (
	"bytes"
	"encoding/binary"
	"michelprogram/photon-parser/photon/parameters/readers"
)

func Parse(reader *bytes.Reader) (*Parameters, error){
	res := &Parameters{}

	var header Header
	if err := binary.Read(reader, binary.BigEndian, &header); err != nil {
		return nil, err
	}
	
	value, err := readers.Decode(reader, header.Type); 
	
	if err != nil {
		return nil, err
	}
	
	res.Header = header
	res.Value = value

	return res,nil
}
