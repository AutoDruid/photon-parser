package assembler

import (
	"bytes"
	"testing"

	"AutoDruid/photon-parser/internal/types"
)

func TestFeedSingleFragmentCompletes(t *testing.T) {
	a := NewAssembler()

	got, complete := a.Feed(types.Fragment{
		ID:     1,
		Size:   3,
		Offset: 0,
		Data:   []byte{0x01, 0x02, 0x03},
	})

	if !complete {
		t.Fatal("Feed() complete = false, want true")
	}

	if !bytes.Equal(got, []byte{0x01, 0x02, 0x03}) {
		t.Errorf("Feed() = %v, want %v", got, []byte{0x01, 0x02, 0x03})
	}

	if len(a.buffers) != 0 {
		t.Errorf("buffers len = %d, want 0", len(a.buffers))
	}
}

func TestFeedReassemblesMultipleFragments(t *testing.T) {
	a := NewAssembler()

	got, complete := a.Feed(types.Fragment{
		ID:     1,
		Count:  2,
		Index:  0,
		Size:   4,
		Offset: 0,
		Data:   []byte{0x01, 0x02},
	})

	if complete {
		t.Fatal("Feed() complete = true, want false")
	}

	if got != nil {
		t.Errorf("Feed() = %v, want nil", got)
	}

	got, complete = a.Feed(types.Fragment{
		ID:     1,
		Count:  2,
		Index:  1,
		Size:   4,
		Offset: 2,
		Data:   []byte{0x03, 0x04},
	})

	if !complete {
		t.Fatal("Feed() complete = false, want true")
	}

	if !bytes.Equal(got, []byte{0x01, 0x02, 0x03, 0x04}) {
		t.Errorf("Feed() = %v, want %v", got, []byte{0x01, 0x02, 0x03, 0x04})
	}

	if len(a.buffers) != 0 {
		t.Errorf("buffers len = %d, want 0", len(a.buffers))
	}
}

func TestFeedReassemblesOutOfOrderFragments(t *testing.T) {
	a := NewAssembler()

	_, complete := a.Feed(types.Fragment{
		ID:     1,
		Count:  2,
		Index:  1,
		Size:   4,
		Offset: 2,
		Data:   []byte{0x03, 0x04},
	})

	if complete {
		t.Fatal("Feed() complete = true, want false")
	}

	got, complete := a.Feed(types.Fragment{
		ID:     1,
		Count:  2,
		Index:  0,
		Size:   4,
		Offset: 0,
		Data:   []byte{0x01, 0x02},
	})

	if !complete {
		t.Fatal("Feed() complete = false, want true")
	}

	if !bytes.Equal(got, []byte{0x01, 0x02, 0x03, 0x04}) {
		t.Errorf("Feed() = %v, want %v", got, []byte{0x01, 0x02, 0x03, 0x04})
	}
}

func TestFeedKeepsIndependentBuffers(t *testing.T) {
	a := NewAssembler()

	_, complete := a.Feed(types.Fragment{
		ID:     1,
		Size:   4,
		Offset: 0,
		Data:   []byte{0x01, 0x02},
	})

	if complete {
		t.Fatal("Feed() complete = true, want false")
	}

	got, complete := a.Feed(types.Fragment{
		ID:     2,
		Size:   2,
		Offset: 0,
		Data:   []byte{0x0a, 0x0b},
	})

	if !complete {
		t.Fatal("Feed() complete = false, want true")
	}

	if !bytes.Equal(got, []byte{0x0a, 0x0b}) {
		t.Errorf("Feed() = %v, want %v", got, []byte{0x0a, 0x0b})
	}

	if _, ok := a.buffers[1]; !ok {
		t.Fatal("buffer for fragment ID 1 was deleted")
	}

	if _, ok := a.buffers[2]; ok {
		t.Fatal("buffer for completed fragment ID 2 still exists")
	}
}
