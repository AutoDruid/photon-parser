package parameters

import (
	"michelprogram/photon-parser/parameters/readers"
	"michelprogram/photon-parser/parser"
)

// Parse reads a complete parameter from the reader.
// Format: Header (1 byte ID + 1 byte Type), followed by the typed value.
//
// The function first reads the parameter header to determine the parameter ID
// and type code, then decodes the value according to that type using the
// Protocol16 decoder.
//
// Returns a Parameters struct containing the ID, Type, and decoded Value,
// or an error if parsing fails.
//
// Example usage:
//
//	param, err := parameters.Parse(reader)
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Parameter %d has value: %v\n", param.ID, param.Value)
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
