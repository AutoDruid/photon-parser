package photon

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Photon struct{}

type InitPacket struct {
	Header   [8]byte
	Reserved [392]byte
}

func ParseInitPacket(packet []byte) (*InitPacket, error) {
	reader := bytes.NewReader(packet)

	var header InitPacket
	err := binary.Read(reader, binary.BigEndian, &header)
	if err != nil {
		return nil, err
	}

	return &header, nil
}

func IsHeader(payload []byte) (bool, error) {

	if len(payload) < 7 {
		return false, errors.New("not long enough to be header")
	}

	//Header as MCRH3110
	return bytes.Equal(payload[0:8], []byte{77, 67, 82, 72, 51, 49, 49, 48}), nil
}
