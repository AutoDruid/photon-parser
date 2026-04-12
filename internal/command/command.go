package command

import (
	"fmt"
	"michelprogram/photon-parser/internal/command/acknowledge"
	"michelprogram/photon-parser/internal/command/ping"
	"michelprogram/photon-parser/internal/command/sendReliable"
	"michelprogram/photon-parser/internal/reader"
)

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

	if header.Length < HEADER_SIZE {
		return fmt.Errorf("command length %d smaller than header size %d", header.Length, HEADER_SIZE)
	}

	c.Header = header
	parsed, err := c.parsePayload(header.Type, r)
	if err != nil {
		rest, _ := r.ReadBytes(int(header.Length - HEADER_SIZE))
		// don't fatal — just store raw for encrypted packets
		c.Payload = UnknownPayload{Raw: rest, Kind: header.Type}
	}else{

	}
	c.Payload = parsed

	return nil
}

func (c Command) parsePayload(t Type, r *reader.Reader) (reader.Payload, error) {
	switch t {
	case SendReliable:
		sd := sendReliable.Reliable{}
		err := sd.Parse(r)
		if err != nil {
			return nil, err
		}
		return sd, nil
	case Ping:
		ping := &ping.Ping{}
		ping.Parse(r)
		return ping, nil
	case Acknowledge:
		acknowledge := &acknowledge.Acknowledge{}
		acknowledge.Parse(r)
		return acknowledge, nil
	default:
		return nil, fmt.Errorf("unknown")
	}
}

func (s *Command) parseHeader(r *reader.Reader) (Header, error) {
	var err error
	var header Header

	b, err := r.ReadUInt8()
	if err != nil {
		return Header{}, err
	}

	header.Type = Type(b)

	header.ChannelID, err = r.ReadUInt8()
	if err != nil {
		return Header{}, err
	}

	header.Flags, err = r.ReadUInt8()
	if err != nil {
		return Header{}, err
	}

	header.ReservedByte, err = r.ReadUInt8()
	if err != nil {
		return Header{}, err
	}

	header.Length, err = r.ReadUInt32()
	if err != nil {
		return Header{}, err
	}

	header.ReliableSequenceNumber, err = r.ReadUInt32()
	if err != nil {
		return Header{}, err
	}

	return header, nil
}

func (c Command) String() string {
	return fmt.Sprintf("Type: %d, ChannelID: %d, Flags: %d, ReservedByte: %d, Length: %d, ReliableSequenceNumber: %d", c.Type, c.ChannelID, c.Flags, c.ReservedByte, c.Length, c.ReliableSequenceNumber)
}
