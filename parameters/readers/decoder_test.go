package readers

import (
	"math"
	"michelprogram/photon-parser/parser"
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		name    string
		ttype   Type
		input   []byte
		want    any
		wantErr bool
	}{
		// Primitives
		{
			name:  "int8 positive",
			ttype: Int8Type,
			input: []byte{0x2A}, // 42
			want:  int8(42),
		},
		{
			name:  "int8 negative",
			ttype: Int8Type,
			input: []byte{0xFF}, // -1
			want:  int8(-1),
		},
		{
			name:  "int16",
			ttype: Int16Type,
			input: []byte{0x03, 0xE8}, // 1000
			want:  int16(1000),
		},
		{
			name:  "int32",
			ttype: Int32Type,
			input: []byte{0x00, 0x00, 0x27, 0x10}, // 10000
			want:  int32(10000),
		},
		{
			name:  "int64",
			ttype: Int64Type,
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0xE8}, // 1000
			want:  int64(1000),
		},
		{
			name:  "float32",
			ttype: Float32Type,
			input: []byte{0x3F, 0x80, 0x00, 0x00}, // 1.0
			want:  float32(1.0),
		},
		{
			name:  "float64",
			ttype: Float64Type,
			input: []byte{0x40, 0x09, 0x21, 0xFB, 0x54, 0x44, 0x2D, 0x18}, // pi
			want:  math.Pi,
		},
		{
			name:  "string",
			ttype: StringType,
			input: []byte{0x00, 0x05, 'h', 'e', 'l', 'l', 'o'},
			want:  "hello",
		},
		{
			name:  "boolean true",
			ttype: BooleanType,
			input: []byte{0x01},
			want:  true,
		},
		{
			name:  "boolean false",
			ttype: BooleanType,
			input: []byte{0x00},
			want:  false,
		},

		// Arrays
		{
			name:  "int8 array",
			ttype: Int8ArrayType,
			input: []byte{
				0x00, 0x00, 0x00, 0x03, // size = 3 (uint32)
				0x01, 0x02, 0x03,
			},
			want: []int8{1, 2, 3},
		},
		{
			name:  "int32 array",
			ttype: Int32ArrayType,
			input: []byte{
				0x00, 0x00, 0x00, 0x02, // size = 2 (uint32)
				0x00, 0x00, 0x00, 0x0A, // 10
				0x00, 0x00, 0x00, 0x14, // 20
			},
			want: []int32{10, 20},
		},
		{
			name:  "string array",
			ttype: StringArrayType,
			input: []byte{
				0x00, 0x00, 0x00, 0x02, // size = 2 (uint32)
				0x00, 0x02, 'h', 'i', // "hi"
				0x00, 0x03, 'b', 'y', 'e', // "bye"
			},
			want: []string{"hi", "bye"},
		},
		{
			name:  "generic array (int32s)",
			ttype: ArrayType,
			input: []byte{
				0x00, 0x02, // size = 2 (uint16)
				byte(Int32Type),        // element type
				0x00, 0x00, 0x00, 0x64, // 100
				0x00, 0x00, 0x00, 0xC8, // 200
			},
			want: []any{int32(100), int32(200)},
		},

		// Error cases
		{
			name:    "int8 truncated",
			ttype:   Int8Type,
			input:   []byte{},
			wantErr: true,
		},
		{
			name:    "int16 truncated",
			ttype:   Int16Type,
			input:   []byte{0x01},
			wantErr: true,
		},
		{
			name:    "int32 truncated",
			ttype:   Int32Type,
			input:   []byte{0x00, 0x00, 0x01},
			wantErr: true,
		},
		{
			name:    "string truncated",
			ttype:   StringType,
			input:   []byte{0x00, 0x05, 'h', 'i'}, // length=5 but only 2 chars
			wantErr: true,
		},
		{
			name:    "array truncated",
			ttype:   Int8ArrayType,
			input:   []byte{0x00, 0x00, 0x00, 0x03, 0x01}, // size=3 (uint32) but only 1 element
			wantErr: true,
		},

		// Unknown type should return empty string and nil error (per current implementation)
		{
			name:  "unknown type",
			ttype: Type(0xFF), // Invalid type
			input: []byte{0x01, 0x02, 0x03},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := parser.NewReader(tt.input)
			got, err := Decode(reader, tt.ttype)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Decode() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Decode() unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() = %v (type: %T), want %v (type: %T)", got, got, tt.want, tt.want)
			}
		})
	}
}

// TestDecodeAllTypes ensures all defined types are handled by Decode
func TestDecodeAllTypes(t *testing.T) {
	// Map of all types to sample valid input
	typeSamples := map[Type][]byte{
		Int8Type:        {0x01},
		Int16Type:       {0x00, 0x01},
		Int32Type:       {0x00, 0x00, 0x00, 0x01},
		Int64Type:       {0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		Float32Type:     {0x00, 0x00, 0x00, 0x00},
		Float64Type:     {0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		StringType:      {0x00, 0x00}, // empty string
		BooleanType:     {0x00},
		Int8ArrayType:   {0x00, 0x00, 0x00, 0x00},                 // empty array (uint32 size)
		Int32ArrayType:  {0x00, 0x00, 0x00, 0x00},                 // empty array (uint32 size)
		StringArrayType: {0x00, 0x00, 0x00, 0x00},                 // empty array (uint32 size)
		ArrayType:       {0x00, 0x00, 0x00, 0x00, byte(Int8Type)}, // empty array with type (uint32 size)
	}

	for ttype, input := range typeSamples {
		t.Run(string(rune(ttype)), func(t *testing.T) {
			reader := parser.NewReader(input)
			_, err := Decode(reader, ttype)
			if err != nil {
				t.Errorf("Decode() failed for type 0x%02x: %v", ttype, err)
			}
		})
	}
}

func TestDecodeReaderPosition(t *testing.T) {
	// Create buffer with multiple values: int8(42), int16(1000)
	input := []byte{
		0x2A,       // int8: 42
		0x03, 0xE8, // int16: 1000
	}
	reader := parser.NewReader(input)

	// Read first value (int8)
	val1, err := Decode(reader, Int8Type)
	if err != nil {
		t.Fatalf("First Decode() failed: %v", err)
	}
	if val1 != int8(42) {
		t.Errorf("First value = %v, want 42", val1)
	}

	// Reader should now be at position 1
	if reader.Len() != 2 {
		t.Errorf("After first read, %d bytes remaining, want 2", reader.Len())
	}

	// Read second value (int16)
	val2, err := Decode(reader, Int16Type)
	if err != nil {
		t.Fatalf("Second Decode() failed: %v", err)
	}
	if val2 != int16(1000) {
		t.Errorf("Second value = %v, want 1000", val2)
	}

	// Reader should now be exhausted
	if reader.Len() != 0 {
		t.Errorf("After second read, %d bytes remaining, want 0", reader.Len())
	}
}

func TestDecodeEmptyReader(t *testing.T) {
	types := []Type{
		Int8Type, Int16Type, Int32Type, Int64Type,
		Float32Type, Float64Type, BooleanType,
	}

	for _, ttype := range types {
		t.Run(string(rune(ttype)), func(t *testing.T) {
			reader := parser.NewReader([]byte{})
			_, err := Decode(reader, ttype)
			if err == nil {
				t.Errorf("Decode() with empty reader should fail for type 0x%02x", ttype)
			}
		})
	}
}

func BenchmarkDecode(b *testing.B) {
	benchmarks := []struct {
		name  string
		ttype Type
		input []byte
	}{
		{
			name:  "int8",
			ttype: Int8Type,
			input: []byte{0x2A},
		},
		{
			name:  "int32",
			ttype: Int32Type,
			input: []byte{0x00, 0x00, 0x27, 0x10},
		},
		{
			name:  "string",
			ttype: StringType,
			input: []byte{0x00, 0x05, 'h', 'e', 'l', 'l', 'o'},
		},
		{
			name:  "int32_array",
			ttype: Int32ArrayType,
			input: []byte{
				0x00, 0x00, 0x00, 0x0A, // 10 elements (uint32)
				0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x02,
				0x00, 0x00, 0x00, 0x03,
				0x00, 0x00, 0x00, 0x04,
				0x00, 0x00, 0x00, 0x05,
				0x00, 0x00, 0x00, 0x06,
				0x00, 0x00, 0x00, 0x07,
				0x00, 0x00, 0x00, 0x08,
				0x00, 0x00, 0x00, 0x09,
				0x00, 0x00, 0x00, 0x0A,
			},
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				reader := parser.NewReader(bm.input)
				_, err := Decode(reader, bm.ttype)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
