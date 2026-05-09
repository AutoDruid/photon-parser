package assembler

import "michelprogram/photon-parser/internal/types"

// PacketBuffer accumulates fragments for one fragmented reliable payload.
type PacketBuffer struct {
	totalSize uint32
	received  uint32
	data      []byte
}

// Assembler reassembles reliable fragmented payloads by fragment ID.
type Assembler struct {
	buffers map[uint32]*PacketBuffer
}

// NewAssembler creates an empty fragment assembler.
func NewAssembler() *Assembler {
	return &Assembler{
		buffers: make(map[uint32]*PacketBuffer),
	}
}

// Feed adds one fragment and returns the full payload when reassembly completes.
func (a *Assembler) Feed(fragment types.Fragment) ([]byte, bool) {
	buf, ok := a.buffers[fragment.ID]
	if !ok {
		buf = &PacketBuffer{
			totalSize: fragment.Size,
			data:      make([]byte, fragment.Size),
		}
		a.buffers[fragment.ID] = buf
	}

	n := copy(buf.data[fragment.Offset:], fragment.Data)
	buf.received += uint32(n)

	if buf.received >= buf.totalSize {
		complete := buf.data
		delete(a.buffers, fragment.ID)
		return complete, true
	}

	return nil, false
}
