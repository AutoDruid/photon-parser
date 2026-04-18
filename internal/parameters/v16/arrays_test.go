package v16_test

import (
	. "michelprogram/photon-parser/internal/parameters/v16"
	"michelprogram/photon-parser/internal/reader"
	"michelprogram/photon-parser/internal/types"
	"reflect"
	"testing"
)

func TestReadInt8Array(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    []int8
		wantErr bool
	}{
		{
			name:  "empty array",
			input: []byte{0x00, 0x78, 0x00, 0x00, 0x00, 0x00}, // size = 0
			want:  []int8{},
		},
		{
			name:  "single element",
			input: []byte{0x00, 0x78, 0x00, 0x00, 0x00, 0x01, 0x2A}, // size=1, value=42
			want:  []int8{42},
		},
		{
			name:  "multiple elements",
			input: []byte{0x00, 0x78, 0x00, 0x00, 0x00, 0x03, 0x01, 0x02, 0x03}, // [1, 2, 3]
			want:  []int8{1, 2, 3},
		},
		{
			name:  "negative values",
			input: []byte{0x00, 0x78, 0x00, 0x00, 0x00, 0x02, 0xFF, 0xFE}, // [-1, -2]
			want:  []int8{-1, -2},
		},
		{
			name:    "truncated size",
			input:   []byte{0x00, 0x78, 0x00, 0x00},
			wantErr: true,
		},
		{
			name:    "truncated data",
			input:   []byte{0x00, 0x78, 0x00, 0x00, 0x00, 0x03, 0x01}, // says 3 elements, only 1 present
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := reader.NewReader(tt.input, reader.Options{
				ParameterParser: &Parameter{},
			})
			p := &Parameter{}
			out := &types.Parameter{}
			err := p.Parse(reader, out)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadInt8Array() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(out.Value, tt.want) {
				t.Errorf("ReadInt8Array() = %v, want %v", out.Value, tt.want)
			}
		})
	}
}

func TestReadInt32Array(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    []int32
		wantErr bool
	}{
		{
			name:  "empty array",
			input: []byte{0x00, 0x6e, 0x00, 0x00, 0x00, 0x00},
			want:  []int32{},
		},
		{
			name: "single element",
			input: []byte{
				0x00, 0x6e,
				0x00, 0x00, 0x00, 0x01, // size = 1
				0x00, 0x00, 0x01, 0x00, // value = 256
			},
			want: []int32{256},
		},
		{
			name: "multiple elements",
			input: []byte{
				0x00, 0x6e,
				0x00, 0x00, 0x00, 0x03, // size = 3
				0x00, 0x00, 0x00, 0x01, // 1
				0x00, 0x00, 0x00, 0x02, // 2
				0x00, 0x00, 0x00, 0x03, // 3
			},
			want: []int32{1, 2, 3},
		},
		{
			name: "negative values",
			input: []byte{
				0x00, 0x6e,
				0x00, 0x00, 0x00, 0x02, // size = 2
				0xFF, 0xFF, 0xFF, 0xFF, // -1
				0xFF, 0xFF, 0xFF, 0xFE, // -2
			},
			want: []int32{-1, -2},
		},
		{
			name:    "truncated data",
			input:   []byte{0x00, 0x6e, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			reader := reader.NewReader(tt.input, reader.Options{
				ParameterParser: &Parameter{},
			})
			p := &Parameter{}
			out := &types.Parameter{}
			err := p.Parse(reader, out)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadInt32Array() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(out.Value, tt.want) {
				t.Errorf("ReadInt32Array() = %v, want %v", out.Value, tt.want)
			}
		})
	}
}

func TestReadStringArray(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    []string
		wantErr bool
	}{
		{
			name:  "empty array",
			input: []byte{0x00, 0x61, 0x00, 0x00, 0x00, 0x00},
			want:  []string{},
		},
		{
			name: "single string",
			input: []byte{
				0x00, 0x61,
				0x00, 0x00, 0x00, 0x01, // size = 1
				0x00, 0x05, 'H', 'e', 'l', 'l', 'o', // "Hello"
			},
			want: []string{"Hello"},
		},
		{
			name: "multiple strings",
			input: []byte{
				0x00, 0x61,
				0x00, 0x00, 0x00, 0x03, // size = 3
				0x00, 0x03, 'f', 'o', 'o', // "foo"
				0x00, 0x03, 'b', 'a', 'r', // "bar"
				0x00, 0x03, 'b', 'a', 'z', // "baz"
			},
			want: []string{"foo", "bar", "baz"},
		},
		{
			name: "empty strings",
			input: []byte{
				0x00, 0x61,
				0x00, 0x00, 0x00, 0x02, // size = 2
				0x00, 0x00, // ""
				0x00, 0x00, // ""
			},
			want: []string{"", ""},
		},
		{
			name: "mixed length strings",
			input: []byte{
				0x00, 0x61,
				0x00, 0x00, 0x00, 0x02, // size = 2
				0x00, 0x01, 'A', // "A"
				0x00, 0x04, 'T', 'e', 's', 't', // "Test"
			},
			want: []string{"A", "Test"},
		},
		{
			name:    "truncated size",
			input:   []byte{0x00, 0x61, 0x00, 0x00},
			wantErr: true,
		},
		{
			name: "truncated string data",
			input: []byte{
				0x00, 0x61,
				0x00, 0x00, 0x00, 0x01, // size = 1
				0x00, 0x05, 'H', 'i', // string length says 5, but only 2 chars
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := reader.NewReader(tt.input, reader.Options{
				ParameterParser: &Parameter{},
			})
			p := &Parameter{}
			out := &types.Parameter{}
			err := p.Parse(reader, out)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadStringArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(out.Value, tt.want) {
				t.Errorf("ReadStringArray() = %v, want %v", p.Value, tt.want)
			}
		})
	}
}

func TestReadArray(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    []any
		wantErr bool
	}{
		{
			name: "empty array",
			input: []byte{
				0x00, 0x79,
				0x00, 0x00, // size = 0
				byte(types.Int8Type), // type (doesn't matter for empty)
			},
			want: []any{},
		},
		{
			name: "int8 array",
			input: []byte{
				0x00, 0x79,
				0x00, 0x03, // size = 3
				byte(types.Int8Type), // type
				0x01, 0x02, 0x03,     // values
			},
			want: []any{int8(1), int8(2), int8(3)},
		},
		{
			name: "nested array",
			input: []byte{
				0x00, 0x79,
				0x00, 0x02, // size = 2 (outer)
				byte(types.ArrayType), // type = array
				0x00, 0x02,            // size = 2 (inner #1)
				byte(types.Int8Type), // type = int8
				0x01, 0x02,           // values
				0x00, 0x03, // size = 3 (inner #2)
				byte(types.Int8Type), // type = int8
				0x03, 0x04, 0x05,     // values
			},
			want: []any{
				[]any{int8(1), int8(2)},
				[]any{int8(3), int8(4), int8(5)},
			},
		},
		{
			name: "int32 array",
			input: []byte{
				0x00, 0x79,
				0x00, 0x02, // size = 2
				byte(types.Int32Type),  // type
				0x00, 0x00, 0x00, 0x0A, // 10
				0x00, 0x00, 0x00, 0x14, // 20
			},
			want: []any{int32(10), int32(20)},
		},
		{
			name: "string array",
			input: []byte{
				0x00, 0x79,
				0x00, 0x02, // size = 2
				byte(types.StringType), // type
				0x00, 0x02, 'H', 'i',   // "Hi"
				0x00, 0x03, 'B', 'y', 'e', // "Bye"
			},
			want: []any{"Hi", "Bye"},
		},
		{
			name: "boolean array",
			input: []byte{
				0x00, 0x79,
				0x00, 0x03, // size = 3
				byte(types.BooleanType), // type
				0x01, 0x00, 0x01,        // true, false, true
			},
			want: []any{true, false, true},
		},
		{
			name:    "truncated size",
			input:   []byte{0x00, 0x79, 0x00},
			wantErr: true,
		},
		{
			name:    "missing type",
			input:   []byte{0x00, 0x79, 0x00, 0x01}, // size but no type
			wantErr: true,
		},
		{
			name: "truncated data",
			input: []byte{
				0x00, 0x79,
				0x00, 0x02, // size = 2
				byte(types.Int8Type),
				0x01, // only 1 element, should be 2
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := reader.NewReader(tt.input, reader.Options{
				ParameterParser: &Parameter{},
			})
			p := &Parameter{}
			out := &types.Parameter{}
			err := p.Parse(reader, out)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(out.Value.([]any), tt.want) {
				t.Errorf("ReadArray() = %v (types: %T), want %v (types: %T)", p.Value, p.Value, tt.want, tt.want)
			}
		})
	}
}

func TestReadArrayGeneric(t *testing.T) {
	t.Run("uint32 array", func(t *testing.T) {
		input := []byte{
			0x00, 0x79,
			0x00, 0x02, // size = 2
			0x69,
			0x00, 0x00, 0x00, 0x64, // 100
			0x00, 0x00, 0x00, 0xC8, // 200
		}
		reader := reader.NewReader(input, reader.Options{
			ParameterParser: &Parameter{},
		})
		p := &Parameter{}
		out := &types.Parameter{}
		err := p.Parse(reader, out)

		if err != nil {
			t.Fatalf("readArray[uint32]() error = %v", err)
		}
		got, ok := out.Value.([]any)
		if !ok {
			t.Fatalf("value type %T, want []any", out.Value)
		}
		wantU32 := []uint32{100, 200}
		if len(got) != len(wantU32) {
			t.Fatalf("len = %d, want %d", len(got), len(wantU32))
		}
		for i := range got {
			v, ok := got[i].(int32)
			if !ok {
				t.Fatalf("elem %d type %T, want int32", i, got[i])
			}
			if uint32(v) != wantU32[i] {
				t.Fatalf("elem %d = %d (uint32 %d), want %d", i, v, uint32(v), wantU32[i])
			}
		}
	})

	t.Run("float32 array", func(t *testing.T) {
		input := []byte{
			0x00, 0x79, // ID, types.ArrayType
			0x00, 0x01, // uint16 size = 1
			0x66,                   // Float32Type
			0x3f, 0x80, 0x00, 0x00, // 1.0 BE
		}
		reader := reader.NewReader(input, reader.Options{
			ParameterParser: &Parameter{},
		})
		p := &Parameter{}
		out := &types.Parameter{}
		err := p.Parse(reader, out)
		if err != nil {
			t.Fatalf("readArray[float32]() error = %v", err)
		}
		got, ok := out.Value.([]any)
		if !ok {
			t.Fatalf("value type %T, want []any", out.Value)
		}
		want := []float32{1.0}
		if len(got) != len(want) {
			t.Fatalf("len = %d, want %d", len(got), len(want))
		}
		for i := range got {
			v, ok := got[i].(float32)
			if !ok {
				t.Fatalf("elem %d type %T, want float32", i, got[i])
			}
			if v != want[i] {
				t.Fatalf("elem %d = %v, want %v", i, v, want[i])
			}
		}
	})
}

func BenchmarkReadInt8Array(b *testing.B) {
	data := []byte{
		0x00, 0x62,
		0x00, 0x00, 0x00, 0x0A, // 10 elements
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A,
	}
	b.ResetTimer()
	reader := reader.NewReader(data, reader.Options{
		ParameterParser: &Parameter{},
	})
	p := &Parameter{}

	for i := 0; i < b.N; i++ {
		out := &types.Parameter{}
		err := p.Parse(reader, out)
		if err != nil {
			b.Fatalf("ReadInt8Array() error = %v", err)
		}
	}
}

func BenchmarkReadStringArray(b *testing.B) {
	data := []byte{
		0x00, 0x61,
		0x00, 0x00, 0x00, 0x02, // 2 elements
		0x00, 0x04, 'T', 'e', 's', 't',
		0x00, 0x04, 'D', 'a', 't', 'a',
	}
	b.ResetTimer()
	reader := reader.NewReader(data, reader.Options{
		ParameterParser: &Parameter{},
	})
	p := &Parameter{}

	for i := 0; i < b.N; i++ {
		out := &types.Parameter{}
		err := p.Parse(reader, out)
		if err != nil {
			b.Fatalf("ReadStringArray() error = %v", err)
		}
	}
}
