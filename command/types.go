package command

type Type uint8

const (
	Acknowledge Type = 0x01
	Connect     Type = 0x02
	VerifyConnect Type = 0x03
	Disconnect    Type = 0x04
	Ping          Type = 0x05
	SendReliable  Type = 0x06
	SendUnreliable Type = 0x07
	SendReliableFragment Type = 0x08
)

const HEADER_SIZE = 12

type Header struct {
	Type                   Type
	ChannelID              uint8
	Flags                  uint8
	ReservedByte           uint8
	Length                 uint32
	ReliableSequenceNumber uint32
}

type Command struct {
	Header

	Data []byte
}
