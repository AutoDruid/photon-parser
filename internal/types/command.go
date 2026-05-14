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
type Command[P ParameterView] struct {
	CommandHeader `json:"header"`

	UnreliablePayload       Reliable[P] `json:"unreliable_payload"`
	ReliablePayload         Reliable[P] `json:"reliable_payload"`
	ReliableFragmentPayload Fragment    `json:"reliable_fragment_payload"`
	AcknowledgePayload      Acknowledge `json:"acknowledge_payload"`
	ConnectPayload          Connect     `json:"connect_payload"`
	UnknownPayload
	PingPayload       struct{} `json:"ping_payload"`
	DisconnectPayload struct{} `json:"disconnect_payload"`
}

type Connect struct {
	Mtu                        uint32 `json:"mtu"`
	WindowSize                 uint32 `json:"window_size"`
	ChannelCount               uint32 `json:"channel_count"`
	IncomingBandwidth          uint32 `json:"incoming_bandwidth"`
	OutgoingBandwidth          uint32 `json:"outgoing_bandwidth"`
	DisconnectThrottle         uint32 `json:"disconnect_throttle"`
	PacketThrottleAcceleration uint32 `json:"packet_throttle_acceleration"`
	PacketThrottleDeceleration uint32 `json:"packet_throttle_deceleration"`
}

type Acknowledge struct {
	AckReliableSequenceNumber uint32 `json:"ack_reliable_sequence_number"`
	AckSentTime               uint32 `json:"ack_sent_time"`
}

type UnknownPayload struct {
	Raw  []byte      `json:"raw"`
	Kind CommandType `json:"kind"`
}

type Fragment struct {
	ID     uint32 `json:"id"`
	Count  uint32 `json:"count"`
	Index  uint32 `json:"index"`
	Size   uint32 `json:"size"`
	Offset uint32 `json:"offset"`
	Data   []byte `json:"data"`
}

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

// ReliableHeader represents the reliable message header.
// This appears at the start of the payload in SendReliable commands.
type ReliableHeader struct {
	Signature      uint8 `json:"signature"`       // Message signature (typically 0xF3)
	Type           Type  `json:"type"`            // Message type (operation, event, etc.)
	EventCode      uint8 `json:"event_code"`      // Operation/event code (application-specific)
	ParameterCount int   `json:"parameter_count"` // Number of parameters following this header
}

type Reliable[P ParameterView] struct {
	ReliableHeader `json:"reliable_header"`
	Parameters     []P `json:"parameters"`
}
