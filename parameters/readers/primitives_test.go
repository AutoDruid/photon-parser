package readers

import (
	"math"
	. "michelprogram/photon-parser/parser"
	"testing"
)

func TestReadInt8(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    int8
		wantErr bool
	}{
		{
			name:  "zero",
			input: []byte{0x00},
			want:  0,
		},
		{
			name:  "positive value",
			input: []byte{0x7F}, // 127
			want:  127,
		},
		{
			name:  "negative value",
			input: []byte{0xFF}, // -1
			want:  -1,
		},
		{
			name:  "negative min",
			input: []byte{0x80}, // -128
			want:  -128,
		},
		{
			name:    "truncated",
			input:   []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewReader(tt.input)
			got, err := ReadInt8(reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadInt8() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("ReadInt8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadInt16(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    int16
		wantErr bool
	}{
		{
			name:  "zero",
			input: []byte{0x00, 0x00},
			want:  0,
		},
		{
			name:  "positive value",
			input: []byte{0x01, 0x00}, // 256 (big endian)
			want:  256,
		},
		{
			name:  "max positive",
			input: []byte{0x7F, 0xFF}, // 32767
			want:  32767,
		},
		{
			name:  "negative value",
			input: []byte{0xFF, 0xFF}, // -1
			want:  -1,
		},
		{
			name:    "truncated single byte",
			input:   []byte{0x00},
			wantErr: true,
		},
		{
			name:    "empty",
			input:   []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewReader(tt.input)
			got, err := ReadInt16(reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadInt16() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("ReadInt16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadInt32(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    int32
		wantErr bool
	}{
		{
			name:  "zero",
			input: []byte{0x00, 0x00, 0x00, 0x00},
			want:  0,
		},
		{
			name:  "small positive",
			input: []byte{0x00, 0x00, 0x01, 0x00}, // 256
			want:  256,
		},
		{
			name:  "large positive",
			input: []byte{0x7F, 0xFF, 0xFF, 0xFF}, // 2147483647
			want:  2147483647,
		},
		{
			name:  "negative",
			input: []byte{0xFF, 0xFF, 0xFF, 0xFF}, // -1
			want:  -1,
		},
		{
			name:    "truncated",
			input:   []byte{0x00, 0x00},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewReader(tt.input)
			got, err := ReadInt32(reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadInt32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("ReadInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadInt64(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    int64
		wantErr bool
	}{
		{
			name:  "zero",
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			want:  0,
		},
		{
			name:  "positive",
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00},
			want:  256,
		},
		{
			name:  "large positive",
			input: []byte{0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
			want:  9223372036854775807, // max int64
		},
		{
			name:  "negative",
			input: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
			want:  -1,
		},
		{
			name:    "truncated",
			input:   []byte{0x00, 0x00, 0x00, 0x00},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewReader(tt.input)
			got, err := ReadInt64(reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("ReadInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadFloat32(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    float32
		wantErr bool
	}{
		{
			name:  "zero",
			input: []byte{0x00, 0x00, 0x00, 0x00},
			want:  0.0,
		},
		{
			name:  "one",
			input: []byte{0x3F, 0x80, 0x00, 0x00}, // 1.0 in IEEE 754
			want:  1.0,
		},
		{
			name:  "negative one",
			input: []byte{0xBF, 0x80, 0x00, 0x00}, // -1.0
			want:  -1.0,
		},
		{
			name:  "pi",
			input: []byte{0x40, 0x49, 0x0F, 0xDB}, // ~3.14159
			want:  3.14159274,
		},
		{
			name:    "truncated",
			input:   []byte{0x00, 0x00},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewReader(tt.input)
			got, err := ReadFloat32(reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFloat32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Use epsilon comparison for floats
				if math.Abs(float64(got-tt.want)) > 0.00001 {
					t.Errorf("ReadFloat32() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestReadFloat64(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    float64
		wantErr bool
	}{
		{
			name:  "zero",
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			want:  0.0,
		},
		{
			name:  "one",
			input: []byte{0x3F, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			want:  1.0,
		},
		{
			name:  "pi",
			input: []byte{0x40, 0x09, 0x21, 0xFB, 0x54, 0x44, 0x2D, 0x18},
			want:  3.141592653589793,
		},
		{
			name:    "truncated",
			input:   []byte{0x00, 0x00, 0x00, 0x00},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewReader(tt.input)
			got, err := ReadFloat64(reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if math.Abs(got-tt.want) > 0.0000001 {
					t.Errorf("ReadFloat64() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestReadString(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    string
		wantErr bool
	}{
		{
			name:  "empty string",
			input: []byte{0x00, 0x00}, // length = 0
			want:  "",
		},
		{
			name:  "short string",
			input: []byte{0x00, 0x05, 'H', 'e', 'l', 'l', 'o'}, // length = 5
			want:  "Hello",
		},
		{
			name:  "single char",
			input: []byte{0x00, 0x01, 'A'}, // length = 1
			want:  "A",
		},
		{
			name:  "unicode string",
			input: []byte{0x00, 0x06, 0xE4, 0xB8, 0xAD, 0xE6, 0x96, 0x87}, // "中文" in UTF-8
			want:  "中文",
		},
		{
			name:    "missing length",
			input:   []byte{},
			wantErr: true,
		},
		{
			name:    "truncated length",
			input:   []byte{0x00},
			wantErr: true,
		},
		{
			name:    "truncated string",
			input:   []byte{0x00, 0x05, 'H', 'i'}, // length says 5, but only 2 bytes
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewReader(tt.input)
			got, err := ReadString(reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("ReadString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReadBoolean(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    bool
		wantErr bool
	}{
		{
			name:  "false",
			input: []byte{0x00},
			want:  false,
		},
		{
			name:  "true",
			input: []byte{0x01},
			want:  true,
		},
		{
			name:    "invalid value 2",
			input:   []byte{0x02},
			wantErr: true,
		},
		{
			name:    "invalid value 255",
			input:   []byte{0xFF},
			wantErr: true,
		},
		{
			name:    "truncated",
			input:   []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewReader(tt.input)
			got, err := ReadBoolean(reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadBoolean() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("ReadBoolean() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadPrimitiveGeneric(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		reader := NewReader([]byte{0xFF})
		got, err := ReadPrimitive[uint8](reader)
		if err != nil {
			t.Fatalf("readPrimitive[uint8]() error = %v", err)
		}
		if got != 255 {
			t.Errorf("readPrimitive[uint8]() = %v, want 255", got)
		}
	})

	t.Run("uint32", func(t *testing.T) {
		reader := NewReader([]byte{0x00, 0x00, 0x01, 0x00})
		got, err := ReadPrimitive[uint32](reader)
		if err != nil {
			t.Fatalf("readPrimitive[uint32]() error = %v", err)
		}
		if got != 256 {
			t.Errorf("readPrimitive[uint32]() = %v, want 256", got)
		}
	})
}

func BenchmarkReadInt32(b *testing.B) {
	data := []byte{0x00, 0x00, 0x01, 0x00}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader := NewReader(data)
		_, _ = ReadInt32(reader)
	}
}

func BenchmarkReadString(b *testing.B) {
	data := []byte{0x00, 0x05, 'H', 'e', 'l', 'l', 'o'}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader := NewReader(data)
		_, _ = ReadString(reader)
	}
}
