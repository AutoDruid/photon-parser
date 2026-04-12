package reliable

import (
	"michelprogram/photon-parser/parameters"
	"michelprogram/photon-parser/parser"
)

type Ping struct{}
func ParsePing() (*Ping, error) {
	return &Ping{}, nil
}

type Acknowledge struct{}
func ParseAcknowledge() (*Acknowledge, error) {
	return &Acknowledge{}, nil
}

// Parse parses a Photon reliable message from a byte slice.
// This is a convenience wrapper around ParseFromReader.
//
// Reliable messages are found in the Data field of SendReliable commands.
// They contain operations, events, and responses with typed parameters.
//
// Returns a Reliable struct containing the parsed header and all parameters,
// or an error if parsing fails.
//
// Example usage:
//
//	// Inside a SendReliable command
//	reliable, err := reliable.Parse(cmd.Data)
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Message type: %d, Event code: %d\n", reliable.Type, reliable.EventCode)
//	for _, param := range reliable.Parameters {
//	    fmt.Printf("  Param %d: %v\n", param.ID, param.Value)
//	}
func Parse(data []byte) (*Reliable, error) {
	return ParseFromReader(parser.NewReader(data))
}

// ParseFromReader parses a Photon reliable message from a parser.Reader.
// It reads the 5-byte header, then iterates through and parses each parameter
// as specified by the ParameterCount field.
//
// The message format consists of:
//   - Header (5 bytes: signature, type, event code, parameter count)
//   - Parameters (ParameterCount entries, each with ID, type, and value)
//
// Returns a Reliable struct with all fields populated including the Parameters slice,
// or an error if any part of parsing fails.
func ParseFromReader(r *parser.Reader) (*Reliable, error) {

	res := &Reliable{}

	header, err := parser.ReadHeader[Header](r)
	if err != nil {
		return nil, err
	}

<<<<<<< Updated upstream
	res.Signature = header.Signature
	res.Type = header.Type
	res.EventCode = header.EventCode
	res.ParameterCount = header.ParameterCount
=======
	if header.Signature != 0xF3 {
        return nil, fmt.Errorf("encrypted or unknown packet, signature: 0x%02x", header.Signature)
    }

	res.Signature = header.Signature
	res.Type = header.Type
	res.EventCode = header.EventCode
	res.ParameterCount = header.ParameterCount
>>>>>>> Stashed changes
	res.Parameters = make([]*parameters.Parameters, header.ParameterCount)

	for i := uint16(0); i < res.ParameterCount; i++ {
		param, err := parameters.Parse(r)
		if err != nil {
			return nil, err
		}
		res.Parameters[i] = param
	}

	return res, nil
}
