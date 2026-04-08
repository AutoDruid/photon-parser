// Package reliable provides parsing for Photon Protocol reliable message layer.
// Reliable messages contain game operations, events, and responses with
// key-value parameters. This layer sits inside SendReliable commands.
package reliable

import "michelprogram/photon-parser/parameters"

// HEADER_SIZE is the size in bytes of a reliable message header (5 bytes).
const HEADER_SIZE = 5

// Type represents a Photon reliable message type.
type Type uint8

// Photon Protocol reliable message types.
// These define the different kinds of reliable messages that can be exchanged.
const (
	OperationRequest       Type = 0x02 // Client requests an operation
	OperationResponse      Type = 0x07 // Server responds to an operation
	OtherOperationResponse Type = 0x03 // Alternative response format
	EventDataType          Type = 0x04 // Server sends an event to client
	ExchangeKeys           Type = 0x06 // Key exchange for encryption
)

// Header represents the reliable message header.
// This appears at the start of the payload in SendReliable commands.
type Header struct {
	Signature      uint8  // Message signature (typically 0xF3)
	Type           Type   // Message type (operation, event, etc.)
	EventCode      uint8  // Operation/event code (application-specific)
	ParameterCount uint16 // Number of parameters following this header
}

// Reliable represents a complete reliable message with header and parameters.
// Parameters contain the actual game data as key-value pairs where each
// parameter has an ID, type, and value.
type Reliable struct {
	Header
	Parameters []*parameters.Parameters // Slice of decoded parameters
}
