package parameters

import (
	"michelprogram/photon-parser/parameters/readers"
	"michelprogram/photon-parser/parser"
)

func Parse(r *parser.Reader) (*Parameters, error) {
	res := &Parameters{}

	header, err := parser.ReadHeader[Header](r)
	if err != nil {
		return nil, err
	}

	value, err := readers.Decode(r, header.Type)

	if err != nil {
		return nil, err
	}

	res.Header = *header
	res.Value = value

	return res, nil
}
