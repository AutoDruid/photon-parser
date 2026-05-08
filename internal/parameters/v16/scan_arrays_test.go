package v16_test

import (
	"reflect"
	"slices"
	"testing"

	v16 "michelprogram/photon-parser/internal/parameters/v16"
	"michelprogram/photon-parser/internal/reader"
)

func TestParseInt8ArrayParameterAndAccessor(t *testing.T) {

	inputs := []struct {
		name  string
		input []byte
		want  []int8
	}{
		{
			name:  "empty array",
			input: []byte{0x01, byte(v16.Int8ArrayType), 0x00, 0x00, 0x00, 0x00}, // size = 0
			want:  []int8{},
		},
		{
			name:  "single element",
			input: []byte{0x01, byte(v16.Int8ArrayType), 0x00, 0x00, 0x00, 0x01, 0x2A}, // size=1, value=42
			want:  []int8{42},
		},
		{
			name:  "multiple elements",
			input: []byte{0x01, byte(v16.Int8ArrayType), 0x00, 0x00, 0x00, 0x03, 0x01, 0x02, 0x03}, // [1, 2, 3]
			want:  []int8{1, 2, 3},
		},
		{
			name:  "negative values",
			input: []byte{0x01, byte(v16.Int8ArrayType), 0x00, 0x00, 0x00, 0x02, 0xFF, 0xFE}, // [-1, -2]
			want:  []int8{-1, -2},
		},
	}

	for _, tt := range inputs {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)
			var parser v16.Parameter
			var got v16.Parameter

			if err := parser.Parse(r, &got, nil); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if got.ID() != 1 {
				t.Errorf("ID() = %d, want 1", got.ID())
			}

			if got.Kind != v16.Int8ArrayType {
				t.Errorf("Kind = %d, want %d", got.Kind, v16.Int8ArrayType)
			}

			if got.Num != uint64(len(tt.want)) {
				t.Errorf("Num = %d, want %d", got.Num, len(tt.want))
			}

			var values []int8
			for _, value := range got.Int8ArrayValue() {
				values = append(values, value)
			}

			if !slices.Equal(values, tt.want) {
				t.Errorf("Int8ArrayValue() = %q, want %q", values, tt.want)
			}

			if r.Cursor != len(tt.input) {
				t.Errorf("Cursor = %d, want %d", r.Cursor, len(tt.input))
			}
		})
	}
}

func TestParseStringArrayParameterAndAccessor(t *testing.T) {

	inputs := []struct {
		name  string
		input []byte
		want  []string
	}{
		{
			name:  "empty array",
			input: []byte{0x01, byte(v16.StringArrayType), 0x00, 0x00, 0x00, 0x00},
			want:  []string{},
		},
		{
			name: "single string",
			input: []byte{
				0x01, byte(v16.StringArrayType),
				0x00, 0x00, 0x00, 0x01, // size = 1
				0x00, 0x05, 'H', 'e', 'l', 'l', 'o', // "Hello"
			},
			want: []string{"Hello"},
		},
		{
			name: "multiple strings",
			input: []byte{
				0x01, byte(v16.StringArrayType),
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
				0x01, byte(v16.StringArrayType),
				0x00, 0x00, 0x00, 0x02, // size = 2
				0x00, 0x00, // ""
				0x00, 0x00, // ""
			},
			want: []string{"", ""},
		},
		{
			name: "mixed length strings",
			input: []byte{
				0x01, byte(v16.StringArrayType),
				0x00, 0x00, 0x00, 0x02, // size = 2
				0x00, 0x01, 'A', // "A"
				0x00, 0x04, 'T', 'e', 's', 't', // "Test"
			},
			want: []string{"A", "Test"},
		},
	}

	for _, tt := range inputs {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)
			var parser v16.Parameter
			var got v16.Parameter

			if err := parser.Parse(r, &got, nil); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if got.ID() != 1 {
				t.Errorf("ID() = %d, want 1", got.ID())
			}

			if got.Kind != v16.StringArrayType {
				t.Errorf("Kind = %d, want %d", got.Kind, v16.StringArrayType)
			}

			if got.Num != uint64(len(tt.want)) {
				t.Errorf("Num = %d, want %d", got.Num, len(tt.want))
			}

			var values []string
			for _, value := range got.StringArrayValue() {
				values = append(values, value)
			}

			if !slices.Equal(values, tt.want) {
				t.Errorf("StringArrayValue() = %q, want %q", values, tt.want)
			}

			if r.Cursor != len(tt.input) {
				t.Errorf("Cursor = %d, want %d", r.Cursor, len(tt.input))
			}
		})
	}
}

func TestParseInt32ArrayParameterAndAccessor(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  []int32
	}{
		{
			name:  "empty array",
			input: []byte{0x01, byte(v16.Int32ArrayType), 0x00, 0x00, 0x00, 0x00},
			want:  []int32{},
		},
		{
			name: "single element",
			input: []byte{
				0x01, byte(v16.Int32ArrayType),
				0x00, 0x00, 0x00, 0x01, // size = 1
				0x00, 0x00, 0x01, 0x00, // value = 256
			},
			want: []int32{256},
		},
		{
			name: "multiple elements",
			input: []byte{
				0x01, byte(v16.Int32ArrayType),
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
				0x01, byte(v16.Int32ArrayType),
				0x00, 0x00, 0x00, 0x02, // size = 2
				0xFF, 0xFF, 0xFF, 0xFF, // -1
				0xFF, 0xFF, 0xFF, 0xFE, // -2
			},
			want: []int32{-1, -2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)
			var parser v16.Parameter
			var got v16.Parameter

			if err := parser.Parse(r, &got, nil); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if got.ID() != 1 {
				t.Errorf("ID() = %d, want 1", got.ID())
			}

			if got.Kind != v16.Int32ArrayType {
				t.Errorf("Kind = %d, want %d", got.Kind, v16.Int32ArrayType)
			}

			if got.Num != uint64(len(tt.want)) {
				t.Errorf("Num = %d, want %d", got.Num, len(tt.want))
			}

			var values []int32
			for _, value := range got.Int32ArrayValue() {
				values = append(values, value)
			}

			if !slices.Equal(values, tt.want) {
				t.Errorf("Int32ArrayValue() = %q, want %q", values, tt.want)
			}

			if r.Cursor != len(tt.input) {
				t.Errorf("Cursor = %d, want %d", r.Cursor, len(tt.input))
			}
		})
	}
}

func TestParseArrayParameterAndAccessor(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  []any
	}{
		{
			name: "empty array",
			input: []byte{
				0x01, byte(v16.ArrayType),
				0x00, 0x00, // size = 0
				byte(v16.Int8Type), // type (doesn't matter for empty)
			},
			want: []any{},
		},
		{
			name: "int8 array",
			input: []byte{
				0x01, byte(v16.ArrayType),
				0x00, 0x03, // size = 3
				byte(v16.Int8Type), // type
				0x01, 0x02, 0x03,   // values
			},
			want: []any{int64(1), int64(2), int64(3)},
		},
		{
			name: "nested array",
			input: []byte{
				0x01, byte(v16.ArrayType),
				0x00, 0x02, // size = 2 (outer)
				byte(v16.ArrayType), // type = array
				0x00, 0x02,          // size = 2 (inner #1)
				byte(v16.Int8Type), // type = int8
				0x01, 0x02,         // values
				0x00, 0x03, // size = 3 (inner #2)
				byte(v16.Int8Type), // type = int8
				0x03, 0x04, 0x05,   // values
			},
			want: []any{
				[]any{int64(1), int64(2)},
				[]any{int64(3), int64(4), int64(5)},
			},
		},
		{
			name: "int32 array",
			input: []byte{
				0x01, byte(v16.ArrayType),
				0x00, 0x02, // size = 2
				byte(v16.Int32Type),    // type
				0x00, 0x00, 0x00, 0x0A, // 10
				0x00, 0x00, 0x00, 0x14, // 20
			},
			want: []any{int64(10), int64(20)},
		},
		{
			name: "string array",
			input: []byte{
				0x01, byte(v16.ArrayType),
				0x00, 0x02, // size = 2
				byte(v16.StringType), // type
				0x00, 0x02, 'H', 'i', // "Hi"
				0x00, 0x03, 'B', 'y', 'e', // "Bye"
			},
			want: []any{"Hi", "Bye"},
		},
		{
			name: "boolean array",
			input: []byte{
				0x01, byte(v16.ArrayType),
				0x00, 0x03, // size = 3
				byte(v16.BooleanType), // type
				0x01, 0x00, 0x01,      // true, false, true
			},
			want: []any{true, false, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)
			var parser v16.Parameter
			var got v16.Parameter

			if err := parser.Parse(r, &got, nil); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if got.ID() != 1 {
				t.Errorf("ID() = %d, want 1", got.ID())
			}

			if got.Kind != v16.ArrayType {
				t.Errorf("Kind = %d, want %d", got.Kind, v16.ArrayType)
			}

			if got.Num != uint64(len(tt.want)) {
				t.Errorf("Num = %d, want %d", got.Num, len(tt.want))
			}

			values := make([]any, 0, got.Num)
			for _, value := range got.ArrayValue() {
				values = append(values, value)
			}

			if !reflect.DeepEqual(values, tt.want) {
				t.Errorf("ArrayValue() = %q, want %q", values, tt.want)
			}

			if r.Cursor != len(tt.input) {
				t.Errorf("Cursor = %d, want %d", r.Cursor, len(tt.input))
			}
		})
	}
}
