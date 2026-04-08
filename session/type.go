// Package session provides parsing for Photon Protocol session layer packets.
// The session layer is the outermost protocol layer, containing session metadata
// and one or more commands.
package session

import "michelprogram/photon-parser/command"

// Header represents the Photon session header containing peer and timing information.
// This header appears at the start of every Photon packet.
type Header struct {
	PeerID       uint16 // Peer identifier for this connection
	CRCEnabled   uint8  // CRC checksum flag (0 = disabled, 1 = enabled)
	CommandCount uint8  // Number of commands following this header
	Timestamp    uint32 // Timestamp in milliseconds
	Challenge    int32  // Challenge value for connection verification
}

// Session represents a complete Photon session packet with its header and commands.
// A session packet can contain multiple commands that will be processed sequentially.
type Session struct {
	Header

	Commands []*command.Command // Slice of commands contained in this session
}
