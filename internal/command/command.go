package command

import (
	"fmt"
	"michelprogram/photon-parser/internal/command/acknowledge"
	"michelprogram/photon-parser/internal/command/ping"
	"michelprogram/photon-parser/internal/command/sendReliable"
	"michelprogram/photon-parser/internal/reader"
	"michelprogram/photon-parser/internal/types"
)

type Command struct {
	types.Command
}

var _ reader.Parseable = (*Command)(nil)

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
func (c *Command) Parse(r *reader.Reader) error {
	header, err := c.parseHeader(r)
	if err != nil {
		return err
	}

	if header.Length < types.COMMAND_HEADER_SIZE {
		return fmt.Errorf(
			"command length %d smaller than header size %d",
			header.Length,
			types.COMMAND_HEADER_SIZE,
		)
	}

	c.CommandHeader = header
	parsed, err := c.parsePayload(header.Type, r)
	if err != nil {
		rest, _ := r.ReadBytes(int(header.Length - types.COMMAND_HEADER_SIZE))
		// don't fatal — just store raw for encrypted packets
		c.Payload = types.UnknownPayload{Raw: rest, Kind: header.Type}
	}

	c.Payload = parsed

	c.emit(r)

	return nil
}

func (c Command) emit(r *reader.Reader) {
	if r.SyncHooks.OnCommand != nil {
		r.SyncHooks.OnCommand(c.Command)
	}
	if r.AsyncHooks.OnCommand != nil {
		select {
		case r.AsyncHooks.OnCommand <- c.Command:
		default: // don't block parser
		}
	}
}
func (c Command) parsePayload(t types.CommandType, r *reader.Reader) (types.Payload, error) {
	switch t {
	case types.SendReliableCommand:
		sd := sendReliable.Reliable{}
		err := sd.Parse(r)
		if err != nil {
			return nil, err
		}
		return sd, nil
	case types.PingCommand:
		ping := &ping.Ping{}
		ping.Parse(r)
		return ping, nil
	case types.AcknowledgeCommand:
		acknowledge := &acknowledge.Acknowledge{}
		acknowledge.Parse(r)
		return acknowledge, nil
	default:
		return nil, fmt.Errorf("unknown")
	}
}

func (s *Command) parseHeader(r *reader.Reader) (types.CommandHeader, error) {
	var err error
	var header types.CommandHeader

	b, err := r.ReadUInt8()
	if err != nil {
		return types.CommandHeader{}, err
	}

	header.Type = types.CommandType(b)

	header.ChannelID, err = r.ReadUInt8()
	if err != nil {
		return types.CommandHeader{}, err
	}

	header.Flags, err = r.ReadUInt8()
	if err != nil {
		return types.CommandHeader{}, err
	}

	header.ReservedByte, err = r.ReadUInt8()
	if err != nil {
		return types.CommandHeader{}, err
	}

	header.Length, err = r.ReadUInt32()
	if err != nil {
		return types.CommandHeader{}, err
	}

	header.ReliableSequenceNumber, err = r.ReadUInt32()
	if err != nil {
		return types.CommandHeader{}, err
	}

	return header, nil
}

func (c Command) String() string {
	return fmt.Sprintf("Type: %d, ChannelID: %d, Flags: %d, ReservedByte: %d, Length: %d, ReliableSequenceNumber: %d", c.Type, c.ChannelID, c.Flags, c.ReservedByte, c.Length, c.ReliableSequenceNumber)
}
