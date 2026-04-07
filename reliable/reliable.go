package reliable

import (
	"michelprogram/photon-parser/parser"
)

func Parse(data []byte) (*Reliable, error) {
    return ParseFromReader(parser.NewReader(data))
}

func ParseFromReader(r *parser.Reader) (*Reliable, error){

	res := &Reliable{}

	header, err := parser.ReadHeader[Header](r)
	if err != nil {
		return nil, err
	}

	res.Signature = header.Signature
	res.Type = header.Type
	res.EventCode = header.EventCode
	res.ParameterCount = header.ParameterCount

	return res,nil
}
