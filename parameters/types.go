// Package parameters provides structures and parsing for Photon Protocol parameters.
// Parameters are key-value pairs used in Photon operations and events, where each
// parameter has an ID, a type code, and a typed value.
package parameters

import (
	"fmt"
	"michelprogram/photon-parser/parameters/readers"
)

// Header represents the parameter header containing the parameter ID and type.
// This appears at the beginning of each serialized parameter.
type Header struct {
	ID   uint8        // Parameter identifier (application-specific)
	Type readers.Type // Protocol16 type code indicating how to decode the value
}

// Parameters represents a complete Photon Protocol parameter with its header and decoded value.
// The Value field contains the decoded data according to the Type specified in the Header.
type Parameters struct {
	Header

	Value interface{} // Decoded value, type depends on Header.Type
}

// String returns a human-readable representation of the parameter.
// Format: "ID: <id>\nType: <type>\nValue: <value>\n"
func (p Parameters) String() string {
	return fmt.Sprintf("ID: %d\nType: %d\nValue: %v\n", p.ID, p.Type, p.Value)
}
