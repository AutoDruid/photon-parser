package wireshark

import (
	"testing"
)

func TestLoadWiresharkJSON(t *testing.T) {
	filepath := "../../ressources/wireshark.json"

	parser, err := LoadFromWiresharkExport(filepath)
	if err != nil {
		t.Fatalf("LoadFromWiresharkExport() failed: %v", err)
	}

	if len(parser.Requests) == 0 {
		t.Fatal("LoadFromWiresharkExport() returned empty requests")
	}

	if len(parser.Requests) < 10 {
		t.Errorf("Expected at least 10 requests, got %d", len(parser.Requests))
	}
}

func TestLoadWiresharkJSON_NonExistentFile(t *testing.T) {
	_, err := LoadFromWiresharkExport("nonexistent.json")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestRequestsDataToBytes(t *testing.T) {
	filepath := "../../ressources/wireshark.json"

	parser, err := LoadFromWiresharkExport(filepath)
	if err != nil {
		t.Fatalf("LoadFromWiresharkExport() failed: %v", err)
	}

	payloads, err := parser.RequestsDataToBytes()
	if err != nil {
		t.Fatalf("RequestsDataToBytes() failed: %v", err)
	}

	if len(payloads) == 0 {
		t.Fatal("RequestsDataToBytes() returned empty slice")
	}

	if len(payloads[0]) == 0 {
		t.Error("First payload is empty")
	}

	firstPayload := payloads[0]

	expectedHeader := []byte{0x4d, 0x43, 0x52, 0x48, 0x33, 0x31, 0x31, 0x30}
	if len(firstPayload) < len(expectedHeader) {
		t.Fatalf("First payload too short: got %d bytes", len(firstPayload))
	}

	for i, b := range expectedHeader {
		if firstPayload[i] != b {
			t.Errorf("Byte %d: expected 0x%02x, got 0x%02x", i, b, firstPayload[i])
		}
	}
}

func TestPayloadHexConversion(t *testing.T) {
	tests := []struct {
		name     string
		requests []WireSharkRequest
		expected [][]byte
		wantErr  bool
	}{
		{
			name: "single request with uppercase hex and separators",
			requests: []WireSharkRequest{
				{
					Source: Source{
						Layers: Layers{
							Data: Data{
								DataData: "4D:43:52:48",
							},
						},
					},
				},
			},
			expected: [][]byte{
				{0x4d, 0x43, 0x52, 0x48},
			},
			wantErr: false,
		},
		{
			name: "multiple requests lowercase and single byte",
			requests: []WireSharkRequest{
				{
					Source: Source{
						Layers: Layers{
							Data: Data{
								DataData: "ff:ff:00:01",
							},
						},
					},
				},
				{
					Source: Source{
						Layers: Layers{
							Data: Data{
								DataData: "2a",
							},
						},
					},
				},
			},
			expected: [][]byte{
				{0xff, 0xff, 0x00, 0x01},
				{0x2a},
			},
			wantErr: false,
		},
		{
			name: "empty string payload",
			requests: []WireSharkRequest{
				{
					Source: Source{
						Layers: Layers{
							Data: Data{
								DataData: "",
							},
						},
					},
				},
			},
			expected: [][]byte{
				{},
			},
			wantErr: false,
		},
		{
			name: "invalid hex should error",
			requests: []WireSharkRequest{
				{
					Source: Source{
						Layers: Layers{
							Data: Data{
								DataData: "zz:01",
							},
						},
					},
				},
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &WiresharkParser{
				Requests: tt.requests,
			}

			got, err := parser.RequestsDataToBytes()
			if (err != nil) != tt.wantErr {
				t.Fatalf("RequestsDataToBytes() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if len(got) != len(tt.expected) {
				t.Fatalf("RequestsDataToBytes() returned %d payloads, expected %d", len(got), len(tt.expected))
			}

			for i := range tt.expected {
				if !bytesEqual(got[i], tt.expected[i]) {
					t.Errorf("payload[%d] = %v, expected %v", i, got[i], tt.expected[i])
				}
			}
		})
	}
}

func TestAllPayloadsAreValid(t *testing.T) {
	filepath := "../../ressources/wireshark.json"
	parser, err := LoadFromWiresharkExport(filepath)
	if err != nil {
		t.Fatalf("LoadFromWiresharkExport() failed: %v", err)
	}

	payloads, err := parser.RequestsDataToBytes()
	if err != nil {
		t.Fatalf("RequestsDataToBytes() failed: %v", err)
	}

	for i, payload := range payloads {
		if len(payload) == 0 {
			t.Errorf("Payload %d is empty", i)
		}

		// All payloads should have at least a few bytes
		if len(payload) < 8 {
			t.Errorf("Payload %d is suspiciously short: %d bytes", i, len(payload))
		}
	}
}

// Helper function to compare byte slices
func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
