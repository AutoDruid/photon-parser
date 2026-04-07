package reliable

import (
	"bytes"
	"encoding/binary"
)

func Parse(reader *bytes.Reader) (*Reliable, error){

	res := &Reliable{}

	var header Header
	if err := binary.Read(reader, binary.BigEndian, &header); err != nil {
		return nil, err
	}

	res.Signature = header.Signature
	res.Type = header.Type
	res.EventCode = header.EventCode
	res.ParameterCount = header.ParameterCount

	return res,nil
}
