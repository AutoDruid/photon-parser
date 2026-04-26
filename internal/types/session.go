package types

// Header represents the Photon session header containing peer and timing information.
// This header appears at the start of every Photon packet.
type Header struct {
	PeerID       uint16 `json:"peer_id"`       // Peer identifier for this connection
	CRCEnabled   uint8  `json:"crc_enabled"`   // CRC checksum flag (0 = disabled, 1 = enabled)
	CommandCount uint8  `json:"command_count"` // Number of commands following this header
	Timestamp    uint32 `json:"timestamp"`     // Timestamp in milliseconds
	Challenge    int32  `json:"challenge"`     // Challenge value for connection verification
}

// Session represents a complete Photon session packet with its header and commands.
// A session packet can contain multiple commands that will be processed sequentially.
type Session struct {
	Header   `json:"header"`
	Commands []Command `json:"commands"`
}
