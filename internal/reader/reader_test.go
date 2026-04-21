package reader_test

import (
	"encoding/binary"
	"math"
	v16 "michelprogram/photon-parser/internal/parameters/v16"
	v18 "michelprogram/photon-parser/internal/parameters/v18"
	"michelprogram/photon-parser/internal/reader"
	"testing"
)

func TestReadByte(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    byte
		wantErr bool
	}{
		{name: "zero", input: []byte{0x00}, want: 0},
		{name: "high bit set is not negative as byte", input: []byte{0xFF}, want: 0xFF},
		{name: "truncated", input: []byte{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input, reader.Options{
				ParameterParser:              &v16.Parameter{},
				ReliableHeaderParameterCount: &v16.ReliableHeaderParameterCountV16{},
				BinaryOrder:                  binary.BigEndian,
			})
			got, err := r.ReadByte()
			if (err != nil) != tt.wantErr {
				t.Fatalf("ReadByte() err = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ReadByte() = %#02x, want %#02x", got, tt.want)
			}
			// If you care about signed interpretation of that same byte:
			if !tt.wantErr && tt.want == 0xFF && int8(got) != -1 {
				t.Errorf("int8(ReadByte()) = %d, want -1", int8(got))
			}
		})
	}
}

func TestReadVarintInt32(t *testing.T) {
	tests := []struct {
		name  string
		want  int32
		input []byte
	}{
		// ── 1-byte varint (zigzag 0–127) ────────────────────────
		{name: "1byte_positive", input: []byte{0x02}, want: 1},
		{name: "1byte_negative", input: []byte{0x01}, want: -1},

		// ── 2-byte varint (zigzag 128–16383) ────────────────────
		{name: "2byte_positive", input: []byte{0x80, 0x01}, want: 64},
		{name: "2byte_negative", input: []byte{0x81, 0x01}, want: -65},

		// ── 3-byte varint (zigzag 16384–2097151) ────────────────
		{name: "3byte_positive", input: []byte{0x80, 0x80, 0x01}, want: 8192},
		{name: "3byte_negative", input: []byte{0x81, 0x80, 0x01}, want: -8193},

		// ── 4-byte varint (zigzag 2097152–268435455) ────────────
		{name: "4byte_positive", input: []byte{0x80, 0x80, 0x80, 0x01}, want: 1048576},
		{name: "4byte_negative", input: []byte{0x81, 0x80, 0x80, 0x01}, want: -1048577},

		// ── 5-byte varint (zigzag 268435456–4294967295) ─────────
		{name: "5byte_positive", input: []byte{0xFE, 0xFF, 0xFF, 0xFF, 0x0F}, want: 2147483647},  // int32 max
		{name: "5byte_negative", input: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x0F}, want: -2147483648}, // int32 min
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input, reader.Options{
				ParameterParser:              &v18.Parameter{},
				ReliableHeaderParameterCount: &v18.ReliableHeaderParameterCountV18{},
				BinaryOrder:                  binary.LittleEndian,
			})
			got, err := r.ReadVarintInt32()
			if err != nil {
				t.Fatalf("ReadVarintInt32() err = %v", err)
			}
			if got != tt.want {
				t.Errorf("ReadVarintInt32() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestReadVarintInt64(t *testing.T) {
	tests := []struct {
		name  string
		want  int64
		input []byte
	}{
		// ── 1-byte varint (zigzag 0–127) ────────────────────────
		{name: "1byte_positive", input: []byte{0x02}, want: 1},
		{name: "1byte_negative", input: []byte{0x01}, want: -1},

		// ── 2-byte varint (zigzag 128–16383) ────────────────────
		{name: "2byte_positive", input: []byte{0x80, 0x01}, want: 64},
		{name: "2byte_negative", input: []byte{0x81, 0x01}, want: -65},

		// ── 3-byte varint (zigzag 16384–2097151) ────────────────
		{name: "3byte_positive", input: []byte{0x80, 0x80, 0x01}, want: 8192},
		{name: "3byte_negative", input: []byte{0x81, 0x80, 0x01}, want: -8193},

		// ── 4-byte varint (zigzag 2097152–268435455) ────────────
		{name: "4byte_positive", input: []byte{0x80, 0x80, 0x80, 0x01}, want: 1048576},
		{name: "4byte_negative", input: []byte{0x81, 0x80, 0x80, 0x01}, want: -1048577},

		// ── 5–9 byte varint: wire value is ZigZag(n), so positive n uses
		// trailing 0x02 (2<<(...)) not 0x01; negative uses 0x81 then 0x80…0x02.
		// ── 5-byte varint ────────────────────────────────────────
		{name: "5byte_positive", input: []byte{0x80, 0x80, 0x80, 0x80, 0x02}, want: 268435456},
		{name: "5byte_negative", input: []byte{0x81, 0x80, 0x80, 0x80, 0x02}, want: -268435457},

		// ── 6-byte varint ────────────────────────────────────────
		{name: "6byte_positive", input: []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x02}, want: 34359738368},
		{name: "6byte_negative", input: []byte{0x81, 0x80, 0x80, 0x80, 0x80, 0x02}, want: -34359738369},

		// ── 7-byte varint ────────────────────────────────────────
		{name: "7byte_positive", input: []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x02}, want: 4398046511104},
		{name: "7byte_negative", input: []byte{0x81, 0x80, 0x80, 0x80, 0x80, 0x80, 0x02}, want: -4398046511105},

		// ── 8-byte varint ────────────────────────────────────────
		{name: "8byte_positive", input: []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x02}, want: 562949953421312},
		{name: "8byte_negative", input: []byte{0x81, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x02}, want: -562949953421313},

		// ── 9-byte varint ────────────────────────────────────────
		{name: "9byte_positive", input: []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x02}, want: 72057594037927936},
		{name: "9byte_negative", input: []byte{0x81, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x02}, want: -72057594037927937},

		// ── 10-byte varint (int64 max/min) ───────────────────────
		{name: "10byte_positive", input: []byte{0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x01}, want: 9223372036854775807},  // int64 max
		{name: "10byte_negative", input: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x01}, want: -9223372036854775808}, // int64 min
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input, reader.Options{
				ParameterParser:              &v18.Parameter{},
				ReliableHeaderParameterCount: &v18.ReliableHeaderParameterCountV18{},
				BinaryOrder:                  binary.LittleEndian,
			})
			got, err := r.ReadVarintInt64()
			if err != nil {
				t.Fatalf("ReadVarintInt64() err = %v", err)
			}
			if got != tt.want {
				t.Errorf("ReadVarintInt64() = %d, want %d", got, tt.want)
			}
		})
	}
}

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
			r := reader.NewReader(tt.input, reader.Options{
				ParameterParser: &v16.Parameter{},
			})
			got, err := r.ReadInt8()

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

func TestReadUInt8(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    uint8
		wantErr bool
	}{
		{name: "zero", input: []byte{0x00}, want: 0},
		{name: "max", input: []byte{0xFF}, want: 255},
		{name: "mid", input: []byte{0x42}, want: 0x42},
		{name: "truncated", input: []byte{}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input, reader.Options{
				ParameterParser:              &v16.Parameter{},
				ReliableHeaderParameterCount: &v16.ReliableHeaderParameterCountV16{},
				BinaryOrder:                  binary.BigEndian,
			})
			got, err := r.ReadUInt8()

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadUInt8() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ReadUInt8() = %v, want %v", got, tt.want)
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
			r := reader.NewReader(tt.input, reader.Options{
				ParameterParser:              &v16.Parameter{},
				ReliableHeaderParameterCount: &v16.ReliableHeaderParameterCountV16{},
				BinaryOrder:                  binary.BigEndian,
			})
			got, err := r.ReadInt16()

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

func TestReadUInt16(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    uint16
		wantErr bool
	}{
		{name: "zero", input: []byte{0x00, 0x00}, want: 0},
		{name: "256", input: []byte{0x01, 0x00}, want: 256},
		{name: "max", input: []byte{0xFF, 0xFF}, want: 65535},
		{name: "truncated single byte", input: []byte{0x00}, wantErr: true},
		{name: "empty", input: []byte{}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input, reader.Options{
				ParameterParser:              &v16.Parameter{},
				ReliableHeaderParameterCount: &v16.ReliableHeaderParameterCountV16{},
				BinaryOrder:                  binary.BigEndian,
			})
			got, err := r.ReadUInt16()

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadUInt16() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ReadUInt16() = %v, want %v", got, tt.want)
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
			r := reader.NewReader(tt.input, reader.Options{
				ParameterParser:              &v16.Parameter{},
				ReliableHeaderParameterCount: &v16.ReliableHeaderParameterCountV16{},
				BinaryOrder:                  binary.BigEndian,
			})
			got, err := r.ReadInt32()

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

func TestReadUInt32(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    uint32
		wantErr bool
	}{
		{name: "zero", input: []byte{0x00, 0x00, 0x00, 0x00}, want: 0},
		{name: "256", input: []byte{0x00, 0x00, 0x01, 0x00}, want: 256},
		{name: "max", input: []byte{0xFF, 0xFF, 0xFF, 0xFF}, want: 4294967295},
		{name: "truncated", input: []byte{0x00, 0x00}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input, reader.Options{
				ParameterParser:              &v16.Parameter{},
				ReliableHeaderParameterCount: &v16.ReliableHeaderParameterCountV16{},
				BinaryOrder:                  binary.BigEndian,
			})
			got, err := r.ReadUInt32()

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadUInt32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ReadUInt32() = %v, want %v", got, tt.want)
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
			r := reader.NewReader(tt.input, reader.Options{
				ParameterParser:              &v16.Parameter{},
				ReliableHeaderParameterCount: &v16.ReliableHeaderParameterCountV16{},
				BinaryOrder:                  binary.BigEndian,
			})
			got, err := r.ReadInt64()

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

func TestReadUInt64(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    uint64
		wantErr bool
	}{
		{
			name:  "zero",
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			want:  0,
		},
		{
			name:  "256",
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00},
			want:  256,
		},
		{
			name:  "max",
			input: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
			want:  18446744073709551615,
		},
		{
			name:    "truncated",
			input:   []byte{0x00, 0x00, 0x00, 0x00},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input, reader.Options{
				ParameterParser:              &v16.Parameter{},
				ReliableHeaderParameterCount: &v16.ReliableHeaderParameterCountV16{},
				BinaryOrder:                  binary.BigEndian,
			})
			got, err := r.ReadUInt64()

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadUInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ReadUInt64() = %v, want %v", got, tt.want)
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
			r := reader.NewReader(tt.input, reader.Options{
				ParameterParser:              &v16.Parameter{},
				ReliableHeaderParameterCount: &v16.ReliableHeaderParameterCountV16{},
				BinaryOrder:                  binary.BigEndian,
			})
			got, err := r.ReadFloat32()

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
			r := reader.NewReader(tt.input, reader.Options{
				ParameterParser:              &v16.Parameter{},
				ReliableHeaderParameterCount: &v16.ReliableHeaderParameterCountV16{},
				BinaryOrder:                  binary.BigEndian,
			})
			got, err := r.ReadFloat64()

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
			r := reader.NewReader(tt.input, reader.Options{
				ParameterParser:              &v16.Parameter{},
				ReliableHeaderParameterCount: &v16.ReliableHeaderParameterCountV16{},
				BinaryOrder:                  binary.BigEndian,
			})
			got, err := r.ReadBoolean()

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

func TestReadString(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		size    int
		want    string
		wantErr bool
	}{
		{
			name:  "empty string",
			input: []byte{},
			size:  0,
			want:  "",
		},
		{
			name:  "short string",
			input: []byte{'H', 'e', 'l', 'l', 'o'}, // length = 5
			size:  5,
			want:  "Hello",
		},
		{
			name:  "single char",
			input: []byte{'A'}, // length = 1
			size:  1,
			want:  "A",
		},
		{
			name:  "unicode string",
			input: []byte{0xE4, 0xB8, 0xAD, 0xE6, 0x96, 0x87}, // "中文" in UTF-8
			size:  6,
			want:  "中文",
		},
		{
			name:    "truncated string",
			input:   []byte{'H', 'i'}, // length says 5, but only 2 bytes
			size:    5,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input, reader.Options{
				ParameterParser:              &v16.Parameter{},
				ReliableHeaderParameterCount: &v16.ReliableHeaderParameterCountV16{},
				BinaryOrder:                  binary.BigEndian,
			})
			got, err := r.ReadString(tt.size)

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
