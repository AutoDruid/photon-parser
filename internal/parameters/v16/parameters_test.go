package v16_test

import (
	"math"
	. "michelprogram/photon-parser/internal/parameters/v16"
	"michelprogram/photon-parser/internal/reader"
	"michelprogram/photon-parser/internal/types"
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		name    string
		ttype   types.ParameterType
		input   []byte
		want    any
		wantErr bool
	}{
		// Primitives
		{
			name:  "int8 positive",
			ttype: types.Int8Type,
			input: []byte{0x2A}, // 42
			want:  int8(42),
		},
		{
			name:  "int8 negative",
			ttype: types.Int8Type,
			input: []byte{0xFF}, // -1
			want:  int8(-1),
		},
		{
			name:  "int16",
			ttype: types.Int16Type,
			input: []byte{0x03, 0xE8}, // 1000
			want:  int16(1000),
		},
		{
			name:  "int32",
			ttype: types.Int32Type,
			input: []byte{0x00, 0x00, 0x27, 0x10}, // 10000
			want:  int32(10000),
		},
		{
			name:  "int64",
			ttype: types.Int64Type,
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0xE8}, // 1000
			want:  int64(1000),
		},
		{
			name:  "float32",
			ttype: types.Float32Type,
			input: []byte{0x3F, 0x80, 0x00, 0x00}, // 1.0
			want:  float32(1.0),
		},
		{
			name:  "float64",
			ttype: types.Float64Type,
			input: []byte{0x40, 0x09, 0x21, 0xFB, 0x54, 0x44, 0x2D, 0x18}, // pi
			want:  math.Pi,
		},
		{
			name:  "string",
			ttype: types.StringType,
			input: []byte{0x00, 0x05, 'h', 'e', 'l', 'l', 'o'},
			want:  "hello",
		},
		{
			name:  "boolean true",
			ttype: types.BooleanType,
			input: []byte{0x01},
			want:  true,
		},
		{
			name:  "boolean false",
			ttype: types.BooleanType,
			input: []byte{0x00},
			want:  false,
		},

		// Arrays
		{
			name:  "int8 array",
			ttype: types.Int8ArrayType,
			input: []byte{
				0x00, 0x00, 0x00, 0x03, // size = 3 (uint32)
				0x01, 0x02, 0x03,
			},
			want: []int8{1, 2, 3},
		},
		{
			name:  "int32 array",
			ttype: types.Int32ArrayType,
			input: []byte{
				0x00, 0x00, 0x00, 0x02, // size = 2 (uint32)
				0x00, 0x00, 0x00, 0x0A, // 10
				0x00, 0x00, 0x00, 0x14, // 20
			},
			want: []int32{10, 20},
		},
		{
			name:  "string array",
			ttype: types.StringArrayType,
			input: []byte{
				0x00, 0x00, 0x00, 0x02, // size = 2 (uint32)
				0x00, 0x02, 'h', 'i', // "hi"
				0x00, 0x03, 'b', 'y', 'e', // "bye"
			},
			want: []string{"hi", "bye"},
		},
		{
			name:  "generic array (int32s)",
			ttype: types.ArrayType,
			input: []byte{
				0x00, 0x02, // size = 2 (uint16)
				byte(types.Int32Type),  // element type
				0x00, 0x00, 0x00, 0x64, // 100
				0x00, 0x00, 0x00, 0xC8, // 200
			},
			want: []any{int32(100), int32(200)},
		},

		// Error cases
		{
			name:    "int8 truncated",
			ttype:   types.Int8Type,
			input:   []byte{},
			wantErr: true,
		},
		{
			name:    "int16 truncated",
			ttype:   types.Int16Type,
			input:   []byte{0x01},
			wantErr: true,
		},
		{
			name:    "int32 truncated",
			ttype:   types.Int32Type,
			input:   []byte{0x00, 0x00, 0x01},
			wantErr: true,
		},
		{
			name:    "string truncated",
			ttype:   types.StringType,
			input:   []byte{0x00, 0x05, 'h', 'i'}, // length=5 but only 2 chars
			wantErr: true,
		},
		{
			name:    "array truncated",
			ttype:   types.Int8ArrayType,
			input:   []byte{0x00, 0x00, 0x00, 0x03, 0x01}, // size=3 (uint32) but only 1 element
			wantErr: true,
		},

		// Unknown type should return empty string and nil error (per current implementation)
		{
			name:    "unknown type",
			ttype:   types.ParameterType(0xFF), // Invalid type
			input:   []byte{0x01, 0x02, 0x03},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fullInput := append([]byte{0x00, byte(tt.ttype)}, tt.input...)
			reader := reader.NewReader(fullInput, reader.Options{
				ParameterParser: &Parameter{},
			})
			param := Parameter{}
			err := param.Parse(reader)

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

			if !reflect.DeepEqual(param.Value, tt.want) {
				t.Errorf("Decode() = %v (type: %T), want %v (type: %T)", param.Value, param.Value, tt.want, tt.want)
			}
		})
	}
}

// TestDecodeAllTypes ensures all defined types are handled by Decode
func TestDecodeAllTypes(t *testing.T) {
	// Map of all types to sample valid input
	typeSamples := map[types.ParameterType][]byte{
		types.Int8Type:        {0x01},
		types.Int16Type:       {0x00, 0x01},
		types.Int32Type:       {0x00, 0x00, 0x00, 0x01},
		types.Int64Type:       {0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		types.Float32Type:     {0x00, 0x00, 0x00, 0x00},
		types.Float64Type:     {0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		types.StringType:      {0x00, 0x00}, // empty string
		types.BooleanType:     {0x00},
		types.Int8ArrayType:   {0x00, 0x00, 0x00, 0x00},           // empty array (uint32 size)
		types.Int32ArrayType:  {0x00, 0x00, 0x00, 0x00},           // empty array (uint32 size)
		types.StringArrayType: {0x00, 0x00, 0x00, 0x00},           // empty array (uint32 size)
		types.ArrayType:       {0x00, 0x00, byte(types.Int8Type)}, // uint16 count 0 + element type
	}

	for ttype, input := range typeSamples {
		t.Run(string(rune(ttype)), func(t *testing.T) {

			fullInput := append([]byte{0x00, byte(ttype)}, input...)
			reader := reader.NewReader(fullInput, reader.Options{
				ParameterParser: &Parameter{},
			})
			param := Parameter{}
			err := param.Parse(reader)

			if err != nil {
				t.Errorf("Decode() failed for type 0x%02x: %v", ttype, err)
			}
		})
	}
}

func TestDecodeReaderPosition(t *testing.T) {
	// Create buffer with multiple values: int8(42), int16(1000)
	input := []byte{
		0x00, byte(types.Int8Type), 0x2A,
		0x00, byte(types.Int16Type), 0x03, 0xE8,
	}
	reader := reader.NewReader(input, reader.Options{
		ParameterParser: &Parameter{},
	})
	param := Parameter{}

	// Read first value (int8)
	err := param.Parse(reader)
	if err != nil {
		t.Fatalf("First Decode() failed: %v", err)
	}
	if param.Value != int8(42) {
		t.Errorf("First value = %v, want 42", param.Value)
	}

	if reader.Cursor != 3 {
		t.Errorf("After first read, %d bytes remaining, want 3", reader.Cursor)
	}

	// Read second value (int16)
	err = param.Parse(reader)
	if err != nil {
		t.Fatalf("Second Decode() failed: %v", err)
	}
	if param.Value != int16(1000) {
		t.Errorf("Second value = %v, want 1000", param.Value)
	}
}

func TestDecodeEmptyReader(t *testing.T) {
	types := []types.ParameterType{
		types.Int8Type, types.Int16Type, types.Int32Type, types.Int64Type,
		types.Float32Type, types.Float64Type, types.BooleanType,
	}

	for _, ttype := range types {
		t.Run(string(rune(ttype)), func(t *testing.T) {
			reader := reader.NewReader([]byte{}, reader.Options{
				ParameterParser: &Parameter{},
			})
			param := Parameter{}
			err := param.Parse(reader)
			if err == nil {
				t.Errorf("Decode() with empty reader should fail for type 0x%02x", ttype)
			}
		})
	}
}

func BenchmarkDecode(b *testing.B) {
	benchmarks := []struct {
		name  string
		ttype types.ParameterType
		input []byte
	}{
		{
			name:  "int8",
			ttype: types.Int8Type,
			input: []byte{0x2A},
		},
		{
			name:  "int32",
			ttype: types.Int32Type,
			input: []byte{0x00, 0x00, 0x27, 0x10},
		},
		{
			name:  "string",
			ttype: types.StringType,
			input: []byte{0x00, 0x05, 'h', 'e', 'l', 'l', 'o'},
		},
		{
			name:  "int32_array",
			ttype: types.Int32ArrayType,
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
			fullInput := append([]byte{0x00, byte(bm.ttype)}, bm.input...)
			r := reader.NewReader(fullInput, reader.Options{
				ParameterParser: &Parameter{},
			})
			param := Parameter{}
			for i := 0; i < b.N; i++ {
				err := param.Parse(r)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
