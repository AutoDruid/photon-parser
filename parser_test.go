package photon_test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/AutoDruid/photon-parser"
	v18 "github.com/AutoDruid/photon-parser/internal/parameters/v18"
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

	frames := loadCaptures("./tests/dataset/v18/1.json", t)

	session := photon.Session{}

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
			err := parser.ParsePacketInto(payload, &session)
			if err != nil {
				t.Fatalf("parse (%d bytes): %v", len(payload), err)
			}
			if &session == nil {
				t.Fatalf("session must not be nil")
			}
		})
	}
	t.Logf("done: %d rows", len(frames))
}

func TestParseUDP2(t *testing.T) {

	parser := photon.NewV18()

	frames := loadCaptures("./tests/dataset/v18/2.json", t)

	session := photon.Session{}
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
			err := parser.ParsePacketInto(payload, &session)
			if err != nil {
				t.Fatalf("parse (%d bytes): %v", len(payload), err)
			}
			if &session == nil {
				t.Fatalf("session must not be nil")
			}
		})
	}
	t.Logf("done: %d rows", len(frames))
}

func TestParseUDP3(t *testing.T) {

	parser := photon.NewV18()

	frames := loadCaptures("./tests/dataset/v18/3.json", t)

	session := photon.Session{}
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
			err := parser.ParsePacketInto(payload, &session)
			if err != nil {
				t.Fatalf("parse (%d bytes): %v", len(payload), err)
			}
			if &session == nil {
				t.Fatalf("session must not be nil")
			}
		})
	}
	t.Logf("done: %d rows", len(frames))
}

func BenchmarkParseUDP1(b *testing.B) {

	parser := photon.NewV18()

	frames := loadCapturesB("./tests/dataset/v18/1.json", b)

	payloads := make([][]byte, len(frames))
	for i, frame := range frames {
		payload, _, err := rowToPayload(frame)
		if err != nil {
			b.Fatalf("row %d: decode row: %v", i, err)
		}
		payloads[i] = payload
	}

	session := photon.Session{}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		for i, payload := range payloads {
			if len(payload) == 0 {
				continue
			}
			err := parser.ParsePacketInto(payload, &session)
			if err != nil {
				b.Fatalf("row %d: error: %v", i, err)
			}
		}
	}
}

func TestCountParameters(t *testing.T) {
	parser := photon.NewV18()

	frames := loadCaptures("./tests/dataset/v18/1.json", t)

	payloads := make([][]byte, len(frames))
	for i, frame := range frames {
		payload, _, err := rowToPayload(frame)
		if err != nil {
			t.Fatalf("row %d: decode row: %v", i, err)
		}
		payloads[i] = payload
	}

	total := 0
	count := 0

	parser.OnEventData(func(r photon.Reliable[v18.Parameter]) {
		total += r.ParameterCount
		count++
	})

	parser.OnOperationRequest(func(r photon.Reliable[v18.Parameter]) {
		total += r.ParameterCount
		count++
	})

	parser.OnOperationResponse(func(r photon.Reliable[v18.Parameter]) {
		total += r.ParameterCount
		count++
	})

	parser.OnOtherOperationResponse(func(r photon.Reliable[v18.Parameter]) {
		total += r.ParameterCount
		count++
	})

	session := photon.Session{}
	for i, payload := range payloads {
		if len(payload) == 0 {
			continue
		}
		err := parser.ParsePacketInto(payload, &session)
		if err != nil {
			t.Fatalf("row %d: error: %v", i, err)
		}
		if &session == nil {
			t.Fatal("sink must not be nil")
		}
	}

	log.Printf("Total %d, avg %.2f", total, float64(total)/float64(count))

}
