package photon_test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"michelprogram/photon-parser"
	"michelprogram/photon-parser/internal/types"
	"os"
	"strings"
	"testing"
)

type WiresharkFrame struct {
	Source struct {
		Layers struct {
			Frame struct {
				Number string `json:"frame.number"`
			} `json:"frame"`
			UDP struct {
				Payload string `json:"udp.payload"`
			} `json:"udp"`
		} `json:"layers"`
	} `json:"_source"`
}

func rowToPayload(item json.RawMessage) (payload []byte, frameSuffix string, err error) {
	var asString string
	if err := json.Unmarshal(item, &asString); err == nil {
		b, err := hex.DecodeString(strings.ReplaceAll(asString, ":", ""))
		return b, "", err
	}

	var row WiresharkFrame
	if err := json.Unmarshal(item, &row); err != nil {
		return nil, "", err
	}
	if n := row.Source.Layers.Frame.Number; n != "" {
		frameSuffix = " [wireshark frame " + n + "]"
	}
	p := row.Source.Layers.UDP.Payload
	if p == "" {
		return nil, frameSuffix, nil
	}
	b, err := hex.DecodeString(strings.ReplaceAll(p, ":", ""))
	return b, frameSuffix, err
}

func loadCaptures(path string, t *testing.T) []json.RawMessage {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Skipf("test data not available: %v", err)
	}
	var frames []json.RawMessage
	if err := json.Unmarshal(data, &frames); err != nil {
		t.Fatalf("invalid test data: %v", err)
	}
	return frames
}

func loadCapturesB(path string, b *testing.B) []json.RawMessage {
	b.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		b.Skipf("test data not available: %v", err)
	}
	var frames []json.RawMessage
	if err := json.Unmarshal(data, &frames); err != nil {
		b.Fatalf("invalid test data: %v", err)
	}
	return frames
}

func TestParseUDP1(t *testing.T) {

	parser := photon.NewV18()

	frames := loadCaptures("resources/v18.json", t)

	for i, frame := range frames {
		payload, frameNo, err := rowToPayload(frame)
		if err != nil {
			t.Fatalf("row %d: decode row: %v", i, err)
		}
		if len(payload) == 0 {
			t.Logf("row %d%s: skip (empty payload)", i, frameNo)
			continue
		}

		name := fmt.Sprintf("row_%d", i)
		if s := strings.TrimSpace(frameNo); s != "" {
			name += "_" + strings.ReplaceAll(strings.ReplaceAll(s, " ", "_"), "__", "_")
		}

		t.Run(name, func(t *testing.T) {
			session, err := parser.ParsePacket(payload)
			if err != nil {
				t.Fatalf("parse (%d bytes): %v", len(payload), err)
			}
			if session == nil {
				t.Fatalf("session must not be nil")
			}
		})
	}
	t.Logf("done: %d rows", len(frames))
}

func TestParseUDP2(t *testing.T) {

	parser := photon.NewV18()

	frames := loadCaptures("resources/v181.json", t)

	for i, frame := range frames {
		payload, frameNo, err := rowToPayload(frame)
		if err != nil {
			t.Fatalf("row %d: decode row: %v", i, err)
		}
		if len(payload) == 0 {
			t.Logf("row %d%s: skip (empty payload)", i, frameNo)
			continue
		}

		name := fmt.Sprintf("row_%d", i)
		if s := strings.TrimSpace(frameNo); s != "" {
			name += "_" + strings.ReplaceAll(strings.ReplaceAll(s, " ", "_"), "__", "_")
		}

		t.Run(name, func(t *testing.T) {
			session, err := parser.ParsePacket(payload)
			if err != nil {
				t.Fatalf("parse (%d bytes): %v", len(payload), err)
			}
			if session == nil {
				t.Fatalf("session must not be nil")
			}
		})
	}
	t.Logf("done: %d rows", len(frames))
}

func TestParseUDP3(t *testing.T) {

	parser := photon.NewV18()

	frames := loadCaptures("resources/v182.json", t)

	for i, frame := range frames {
		payload, frameNo, err := rowToPayload(frame)
		if err != nil {
			t.Fatalf("row %d: decode row: %v", i, err)
		}
		if len(payload) == 0 {
			t.Logf("row %d%s: skip (empty payload)", i, frameNo)
			continue
		}

		name := fmt.Sprintf("row_%d", i)
		if s := strings.TrimSpace(frameNo); s != "" {
			name += "_" + strings.ReplaceAll(strings.ReplaceAll(s, " ", "_"), "__", "_")
		}

		t.Run(name, func(t *testing.T) {
			session, err := parser.ParsePacket(payload)
			if err != nil {
				t.Fatalf("parse (%d bytes): %v", len(payload), err)
			}
			if session == nil {
				t.Fatalf("session must not be nil")
			}
		})
	}
	t.Logf("done: %d rows", len(frames))
}

func BenchmarkParseUDP1(b *testing.B) {

	parser := photon.NewV18()

	frames := loadCapturesB("ressources/v18.json", b)

	b.ReportAllocs()
	b.ResetTimer()

	payloads := make([][]byte, len(frames))
	for i, frame := range frames {
		payload, _, err := rowToPayload(frame)
		if err != nil {
			b.Fatalf("row %d: decode row: %v", i, err)
		}
		payloads[i] = payload
	}

	sess := &types.Session{}

	for b.Loop() {
		for _, payload := range payloads {
			sess, _ = parser.ParsePacket(payload)
		}
	}

	log.Println(sess)
}
