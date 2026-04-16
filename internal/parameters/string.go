package parameters

import "michelprogram/photon-parser/internal/reader"

// ReadString reads a Photon Protocol16 string from the reader.
// Format: uint16 length (big-endian) followed by UTF-8 bytes.
// Returns an empty string if length is 0.
// Returns an error if the declared length cannot be fully read.
//
// Example wire format for "hello":
//
//	0x00 0x05 'h' 'e' 'l' 'l' 'o'
func (p Parameter) readString(reader *reader.Reader) (string, error) {
	size, err := reader.ReadInt16()
	if err != nil {
		return "", err
	}

	return reader.ReadString(int(size))
}
