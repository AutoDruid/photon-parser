package types

// Type represents a Photon Protocol command type.
type CommandType uint8

// Photon Protocol command types.
// These define the various operations that can be performed in a Photon session.
const (
	AcknowledgeCommand          CommandType = 0x01 // Acknowledges receipt of reliable commands
	ConnectCommand              CommandType = 0x02 // Initiates a connection
	VerifyConnectCommand        CommandType = 0x03 // Verifies connection establishment
	DisconnectCommand           CommandType = 0x04 // Gracefully closes a connection
	PingCommand                 CommandType = 0x05 // Keep-alive ping message
	SendReliableCommand         CommandType = 0x06 // Sends reliable data (guaranteed delivery)
	SendUnreliableCommand       CommandType = 0x07 // Sends unreliable data (best effort)
	SendReliableFragmentCommand CommandType = 0x08 // Sends a fragment of a large reliable message
)

// COMMAND_HEADER_SIZE is the size in bytes of a command header (12 bytes).
const COMMAND_HEADER_SIZE = 12

// Header represents the command header containing command metadata.
// This header appears at the start of every command within a session.
type CommandHeader struct {
	Type                   CommandType `json:"type"`                     // Command type (see Type constants)
	ChannelID              uint8       `json:"channel_id"`               // Channel identifier for message ordering
	Flags                  uint8       `json:"flags"`                    // Command flags (implementation-specific)
	ReservedByte           uint8       `json:"reserved_byte"`            // Reserved for future use
	Length                 uint32      `json:"length"`                   // Total length of command including header
	ReliableSequenceNumber uint32      `json:"reliable_sequence_number"` // Sequence number for reliable delivery
}

// Command represents a complete Photon command with its header and payload data.
// The Data field contains the command-specific payload, which may be empty
// for some command types (e.g., Acknowledge, Ping).
type Command struct {
	CommandHeader

	Payload Payload `json:"payload"` // Command payload (interpretation depends on Type)
}

type UnknownPayload struct {
	Raw  []byte
	Kind CommandType
}
