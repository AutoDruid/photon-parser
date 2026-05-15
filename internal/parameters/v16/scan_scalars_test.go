package v16_test

import (
	"bytes"
	"math"
	"testing"

	v16 "github.com/AutoDruid/photon-parser/internal/parameters/v16"
	"github.com/AutoDruid/photon-parser/internal/reader"
)

func TestParseStringParameterAndAccessor(t *testing.T) {
	longString := bytes.Repeat([]byte{'a'}, 128)
	tests := []struct {
		name       string
		input      []byte
		wantID     uint8
		wantString string
		wantNum    uint64
		wantCursor int
	}{
		{
			name:       "string with 1 byte uint16 size",
			input:      append([]byte{0x01, byte(v16.StringType), 0x00, 0x05}, []byte("hello")...),
			wantID:     1,
			wantString: "hello",
			wantNum:    5,
			wantCursor: 9,
		},
		{
			name:       "string with 2 byte uint16 size",
			input:      append([]byte{0x02, byte(v16.StringType), 0x00, 0x80}, longString...),
			wantID:     2,
			wantString: string(longString),
			wantNum:    128,
			wantCursor: 132,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)
			var parser v16.Parameter
			var got v16.Parameter
			if err := parser.ParseInto(r, nil, &got); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if got.ID() != tt.wantID {
				t.Errorf("ID() = %d, want %d", got.ID(), tt.wantID)
			}
			if got.Kind != v16.StringType {
				t.Errorf("Kind = %d, want %d", got.Kind, v16.StringType)
			}
			if got.Num != tt.wantNum {
				t.Errorf("Num = %d, want %d", got.Num, tt.wantNum)
			}
			if gotString, ok := got.StringValue(); !ok || gotString != tt.wantString {
				t.Errorf("StringValue() = %q, want %q", gotString, tt.wantString)
			}
			if r.Cursor != tt.wantCursor {
				t.Errorf("Cursor = %d, want %d", r.Cursor, tt.wantCursor)
			}
		})
	}
}

func TestParseInt8ParameterAndAccessor(t *testing.T) {
	tests := []struct {
		name   string
		input  []byte
		wantID uint8
		want   int8
	}{
		{
			name:   "int8 positive",
			input:  []byte{0x01, byte(v16.Int8Type), 0x2A},
			wantID: 1,
			want:   int8(42),
		},
		{
			name:   "int8 negative",
			input:  []byte{0x02, byte(v16.Int8Type), 0xFF},
			wantID: 2,
			want:   int8(-1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)
			var parser v16.Parameter
			var got v16.Parameter
			if err := parser.ParseInto(r, nil, &got); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if got.ID() != tt.wantID {
				t.Errorf("ID() = %d, want %d", got.ID(), tt.wantID)
			}
			if got.Kind != v16.Int8Type {
				t.Errorf("Kind = %d, want %d", got.Kind, v16.Int8Type)
			}
			if gotInt8, ok := got.IntValue(); !ok || gotInt8 != int64(tt.want) {
				t.Errorf("IntValue() = %d, want %d", gotInt8, tt.want)
			}
		})
	}
}

func TestParseInt16ParameterAndAccessor(t *testing.T) {
	tests := []struct {
		name   string
		input  []byte
		wantID uint8
		want   int16
	}{
		{
			name:   "int16",
			input:  []byte{0x01, byte(v16.Int16Type), 0x03, 0xE8}, // 1000
			wantID: 1,
			want:   int16(1000),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)
			var parser v16.Parameter
			var got v16.Parameter
			if err := parser.ParseInto(r, nil, &got); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if got.ID() != tt.wantID {
				t.Errorf("ID() = %d, want %d", got.ID(), tt.wantID)
			}
			if got.Kind != v16.Int16Type {
				t.Errorf("Kind = %d, want %d", got.Kind, v16.Int16Type)
			}
			if gotInt16, ok := got.IntValue(); !ok || gotInt16 != int64(tt.want) {
				t.Errorf("IntValue() = %d, want %d", gotInt16, tt.want)
			}
		})
	}
}

func TestParseInt32ParameterAndAccessor(t *testing.T) {
	tests := []struct {
		name   string
		input  []byte
		wantID uint8
		want   int32
	}{
		{
			name:   "int32",
			input:  []byte{0x01, byte(v16.Int32Type), 0x00, 0x00, 0x27, 0x10}, // 10000
			wantID: 1,
			want:   int32(10000),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)
			var parser v16.Parameter
			var got v16.Parameter
			if err := parser.ParseInto(r, nil, &got); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if got.ID() != tt.wantID {
				t.Errorf("ID() = %d, want %d", got.ID(), tt.wantID)
			}
			if got.Kind != v16.Int32Type {
				t.Errorf("Kind = %d, want %d", got.Kind, v16.Int32Type)
			}
			if gotInt32, ok := got.IntValue(); !ok || gotInt32 != int64(tt.want) {
				t.Errorf("IntValue() = %d, want %d", gotInt32, tt.want)
			}
		})
	}
}

func TestParseInt64ParameterAndAccessor(t *testing.T) {
	tests := []struct {
		name   string
		input  []byte
		wantID uint8
		want   int64
	}{
		{
			name:   "int64",
			input:  []byte{0x01, byte(v16.Int64Type), 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0xE8}, // 1000
			wantID: 1,
			want:   int64(1000),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)
			var parser v16.Parameter
			var got v16.Parameter
			if err := parser.ParseInto(r, nil, &got); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if got.ID() != tt.wantID {
				t.Errorf("ID() = %d, want %d", got.ID(), tt.wantID)
			}
			if got.Kind != v16.Int64Type {
				t.Errorf("Kind = %d, want %d", got.Kind, v16.Int64Type)
			}
			if gotInt64, ok := got.IntValue(); !ok || gotInt64 != int64(tt.want) {
				t.Errorf("IntValue() = %d, want %d", gotInt64, tt.want)
			}
		})
	}
}

func TestParseFloat32ParameterAndAccessor(t *testing.T) {
	tests := []struct {
		name   string
		input  []byte
		wantID uint8
		want   float32
	}{
		{
			name:   "float32",
			input:  []byte{0x01, byte(v16.Float32Type), 0x3F, 0x80, 0x00, 0x00}, // 1.0
			wantID: 1,
			want:   float32(1.0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)
			var parser v16.Parameter
			var got v16.Parameter
			if err := parser.ParseInto(r, nil, &got); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if got.ID() != tt.wantID {
				t.Errorf("ID() = %d, want %d", got.ID(), tt.wantID)
			}
			if got.Kind != v16.Float32Type {
				t.Errorf("Kind = %d, want %d", got.Kind, v16.Float32Type)
			}
			if gotFloat32, ok := got.Float32Value(); !ok || gotFloat32 != tt.want {
				t.Errorf("Float32Value() = %f, want %f", gotFloat32, tt.want)
			}
		})
	}
}

func TestParseFloat64ParameterAndAccessor(t *testing.T) {
	tests := []struct {
		name   string
		input  []byte
		wantID uint8
		want   float64
	}{
		{
			name:   "float64",
			input:  []byte{0x01, byte(v16.Float64Type), 0x40, 0x09, 0x21, 0xFB, 0x54, 0x44, 0x2D, 0x18}, // pi
			wantID: 1,
			want:   math.Pi,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)
			var parser v16.Parameter
			var got v16.Parameter
			if err := parser.ParseInto(r, nil, &got); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if got.ID() != tt.wantID {
				t.Errorf("ID() = %d, want %d", got.ID(), tt.wantID)
			}
			if got.Kind != v16.Float64Type {
				t.Errorf("Kind = %d, want %d", got.Kind, v16.Float64Type)
			}
			if gotFloat64, ok := got.Float64Value(); !ok || gotFloat64 != tt.want {
				t.Errorf("Float64Value() = %f, want %f", gotFloat64, tt.want)
			}
		})
	}
}

func TestParseBooleanParameterAndAccessor(t *testing.T) {
	tests := []struct {
		name   string
		input  []byte
		wantID uint8
		want   bool
	}{
		{
			name:   "boolean true",
			input:  []byte{0x01, byte(v16.BooleanType), 0x01},
			wantID: 1,
			want:   true,
		},
		{
			name:   "boolean false",
			input:  []byte{0x02, byte(v16.BooleanType), 0x00},
			wantID: 2,
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)
			var parser v16.Parameter
			var got v16.Parameter
			if err := parser.ParseInto(r, nil, &got); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if got.ID() != tt.wantID {
				t.Errorf("ID() = %d, want %d", got.ID(), tt.wantID)
			}
			if got.Kind != v16.BooleanType {
				t.Errorf("Kind = %d, want %d", got.Kind, v16.BooleanType)
			}
			if gotBoolean, ok := got.BooleanValue(); !ok || gotBoolean != tt.want {
				t.Errorf("BooleanValue() = %t, want %t", gotBoolean, tt.want)
			}
		})
	}
}
