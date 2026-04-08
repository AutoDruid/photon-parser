package command

import (
	"fmt"
	"io"
	"michelprogram/photon-parser/parser"
	"sync"
)

// payloadBufPool holds reusable byte slices used as read buffers for command
// payloads. After reading, the bytes are copied into a permanent slice so the
// pooled buffer can be safely returned immediately.
var payloadBufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 512)
		return &b
	},
}

// Parse parses a Photon command from a byte slice.
// This is a convenience wrapper around ParseFromReader.
//
// The command format consists of:
//   - Command header (12 bytes: type, channel, flags, reserved, length, sequence number)
//   - Command payload (length - 12 bytes)
//
// Returns a Command struct containing the parsed header and payload data,
// or an error if parsing fails.
//
// Example usage:
//
//	cmd, err := command.Parse(commandBytes)
//	if err != nil {
//	    return err
//	}
//	if cmd.Type == command.SendReliable {
//	    // Process reliable message
//	}
func Parse(data []byte) (*Command, error) {
	r := parser.NewReaderFromPool(data)
	defer parser.ReleaseReader(r)
	return ParseFromReader(r)
}

// ParseFromReader parses a Photon command from a parser.Reader.
// It first reads the 12-byte command header, validates the length field,
// then reads the remaining payload data.
//
// Returns an error if:
//   - The header cannot be read
//   - The length field is smaller than the header size (invalid)
//   - The payload data cannot be fully read
//
// The returned Command struct contains all header fields and the raw payload
// in the Data field. For SendReliable commands, the Data can be further parsed
// using the reliable package.
func ParseFromReader(r *parser.Reader) (*Command, error) {
	header, err := parser.ReadHeader[Header](r)
	if err != nil {
		return nil, err
	}

	if header.Length < HEADER_SIZE {
		return nil, fmt.Errorf("command length %d smaller than header size %d", header.Length, HEADER_SIZE)
	}

	payloadLen := int(header.Length - HEADER_SIZE)

	// Acquire a pooled buffer, grow if needed, then read the payload.
	bufp := payloadBufPool.Get().(*[]byte)
	buf := *bufp
	if cap(buf) < payloadLen {
		buf = make([]byte, payloadLen)
	}
	buf = buf[:payloadLen]

	if _, err := io.ReadFull(r, buf); err != nil {
		payloadBufPool.Put(bufp)
		return nil, fmt.Errorf("failed to read %d bytes: %w", payloadLen, err)
	}

	// Copy into a permanent slice so cmd.Data outlives the pooled buffer.
	payload := make([]byte, payloadLen)
	copy(payload, buf)
	*bufp = buf
	payloadBufPool.Put(bufp)

	cmd := &Command{}
	cmd.Header = *header
	cmd.Data = payload

	return cmd, nil
}
