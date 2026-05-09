package v18_test

import (
	"reflect"
	"slices"
	"testing"

	v18 "github.com/AutoDruid/photon-parser/internal/parameters/v18"
	"github.com/AutoDruid/photon-parser/internal/reader"
)

func TestParseStringArrayParameterAndAccessor(t *testing.T) {
	want := []string{
		"R2H_BOW:HEAD_LEATHER_SET2:ARMOR_LEATHER_SET1:SHOES_LEATHER_AVALON:CAPEITEM_SMUGGLER",
		"S2H_BOW:HEAD_PLATE_AVAON:ARMOR_CLOTH_AVALON:SHOES_CLOTH_AVALON:CAPEITEM_FW_MARTLOCK",
		"R2H_BOW:HEAD_PLATE_AVALON:ARMOR_CLOTH_FEY:SHOES_LEATHER_AVALON:CAPEITEM_FW_THETFORD",
	}

	input := []byte{
		0x01, byte(v18.StringArrayType), 0x03,
		0x53, 0x52, 0x32, 0x48, 0x5f, 0x42, 0x4f, 0x57, 0x3a, 0x48, 0x45, 0x41, 0x44, 0x5f, 0x4c, 0x45, 0x41, 0x54, 0x48, 0x45, 0x52, 0x5f, 0x53, 0x45, 0x54, 0x32, 0x3a, 0x41, 0x52, 0x4d, 0x4f, 0x52, 0x5f, 0x4c, 0x45, 0x41, 0x54, 0x48, 0x45, 0x52, 0x5f, 0x53, 0x45, 0x54, 0x31, 0x3a, 0x53, 0x48, 0x4f, 0x45, 0x53, 0x5f, 0x4c, 0x45, 0x41, 0x54, 0x48, 0x45, 0x52, 0x5f, 0x41, 0x56, 0x41, 0x4c, 0x4f, 0x4e, 0x3a, 0x43, 0x41, 0x50, 0x45, 0x49, 0x54, 0x45, 0x4d, 0x5f, 0x53, 0x4d, 0x55, 0x47, 0x47, 0x4c, 0x45, 0x52,
		0x53, 0x53, 0x32, 0x48, 0x5f, 0x42, 0x4f, 0x57, 0x3a, 0x48, 0x45, 0x41, 0x44, 0x5f, 0x50, 0x4c, 0x41, 0x54, 0x45, 0x5f, 0x41, 0x56, 0x41, 0x4f, 0x4e, 0x3a, 0x41, 0x52, 0x4d, 0x4f, 0x52, 0x5f, 0x43, 0x4c, 0x4f, 0x54, 0x48, 0x5f, 0x41, 0x56, 0x41, 0x4c, 0x4f, 0x4e, 0x3a, 0x53, 0x48, 0x4f, 0x45, 0x53, 0x5f, 0x43, 0x4c, 0x4f, 0x54, 0x48, 0x5f, 0x41, 0x56, 0x41, 0x4c, 0x4f, 0x4e, 0x3a, 0x43, 0x41, 0x50, 0x45, 0x49, 0x54, 0x45, 0x4d, 0x5f, 0x46, 0x57, 0x5f, 0x4d, 0x41, 0x52, 0x54, 0x4c, 0x4f, 0x43, 0x4b,
		0x53, 0x52, 0x32, 0x48, 0x5f, 0x42, 0x4f, 0x57, 0x3a, 0x48, 0x45, 0x41, 0x44, 0x5f, 0x50, 0x4c, 0x41, 0x54, 0x45, 0x5f, 0x41, 0x56, 0x41, 0x4c, 0x4f, 0x4e, 0x3a, 0x41, 0x52, 0x4d, 0x4f, 0x52, 0x5f, 0x43, 0x4c, 0x4f, 0x54, 0x48, 0x5f, 0x46, 0x45, 0x59, 0x3a, 0x53, 0x48, 0x4f, 0x45, 0x53, 0x5f, 0x4c, 0x45, 0x41, 0x54, 0x48, 0x45, 0x52, 0x5f, 0x41, 0x56, 0x41, 0x4c, 0x4f, 0x4e, 0x3a, 0x43, 0x41, 0x50, 0x45, 0x49, 0x54, 0x45, 0x4d, 0x5f, 0x46, 0x57, 0x5f, 0x54, 0x48, 0x45, 0x54, 0x46, 0x4f, 0x52, 0x44,
	}

	r := reader.NewReader(input)

	var parser v18.Parameter
	var got v18.Parameter

	if err := parser.Parse(r, &got, nil); err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if got.ID() != 1 {
		t.Errorf("ID() = %d, want 1", got.ID())
	}

	if got.Kind != v18.StringArrayType {
		t.Errorf("Kind = %d, want %d", got.Kind, v18.StringArrayType)
	}

	if got.Num != uint64(len(want)) {
		t.Errorf("Num = %d, want %d", got.Num, len(want))
	}

	var values []string
	for _, value := range got.StringArrayValue() {
		values = append(values, value)
	}

	if !slices.Equal(values, want) {
		t.Errorf("StringArrayValue() = %q, want %q", values, want)
	}

	if r.Cursor != len(input) {
		t.Errorf("Cursor = %d, want %d", r.Cursor, len(input))
	}
}

func TestParseStringArrayRolesParameterAndAccessor(t *testing.T) {
	want := []string{"owner", "friend", "user"}

	input := []byte{
		0x02, byte(v18.StringArrayType), 0x03,
		0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72,
		0x06, 0x66, 0x72, 0x69, 0x65, 0x6e, 0x64,
		0x04, 0x75, 0x73, 0x65, 0x72,
	}

	r := reader.NewReader(input)

	var parser v18.Parameter
	var got v18.Parameter

	if err := parser.Parse(r, &got, nil); err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if got.ID() != 2 {
		t.Errorf("ID() = %d, want 2", got.ID())
	}

	if got.Kind != v18.StringArrayType {
		t.Errorf("Kind = %d, want %d", got.Kind, v18.StringArrayType)
	}

	if got.Num != uint64(len(want)) {
		t.Errorf("Num = %d, want %d", got.Num, len(want))
	}

	var values []string
	for _, value := range got.StringArrayValue() {
		values = append(values, value)
	}

	if !slices.Equal(values, want) {
		t.Errorf("StringArrayValue() = %q, want %q", values, want)
	}

	if r.Cursor != len(input) {
		t.Errorf("Cursor = %d, want %d", r.Cursor, len(input))
	}
}

func TestParseFloat32ArrayParametersAndAccessor(t *testing.T) {
	tests := []struct {
		name  string
		id    uint8
		input []byte
		want  []float32
	}{
		{
			name: "id 2",
			id:   2,
			input: []byte{
				0x02, byte(v18.Float32ArrayType), 0x05,
				0x00, 0x00, 0xc8, 0x42,
				0x5d, 0x7e, 0x8a, 0x43,
				0x5d, 0x7e, 0x8a, 0x43,
				0x00, 0x00, 0xc8, 0x42,
				0x5d, 0x7e, 0x8a, 0x43,
			},
			want: []float32{100, 276.9872, 276.9872, 100, 276.9872},
		},
		{
			name: "id 3",
			id:   3,
			input: []byte{
				0x03, byte(v18.Float32ArrayType), 0x05,
				0x00, 0x00, 0xc8, 0x42,
				0x21, 0xfa, 0x3c, 0x43,
				0x21, 0xfa, 0x3c, 0x43,
				0x00, 0x00, 0xc8, 0x42,
				0x21, 0xfa, 0x3c, 0x43,
			},
			want: []float32{100, 188.97707, 188.97707, 100, 188.97707},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)

			var parser v18.Parameter
			var got v18.Parameter

			if err := parser.Parse(r, &got, nil); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if got.ID() != tt.id {
				t.Errorf("ID() = %d, want %d", got.ID(), tt.id)
			}

			if got.Kind != v18.Float32ArrayType {
				t.Errorf("Kind = %d, want %d", got.Kind, v18.Float32ArrayType)
			}

			if got.Num != uint64(len(tt.want)) {
				t.Errorf("Num = %d, want %d", got.Num, len(tt.want))
			}

			var values []float32
			for _, value := range got.Float32ArrayValue() {
				values = append(values, value)
			}

			if !slices.Equal(values, tt.want) {
				t.Errorf("Float32ArrayValue() = %v, want %v", values, tt.want)
			}

			if r.Cursor != len(tt.input) {
				t.Errorf("Cursor = %d, want %d", r.Cursor, len(tt.input))
			}
		})
	}
}

func TestParseInt16ArrayParametersAndAccessor(t *testing.T) {
	tests := []struct {
		name  string
		id    uint8
		input []byte
		want  []int16
	}{
		{
			name: "id 1 values",
			id:   1,
			input: []byte{
				0x01, byte(v18.ShortArrayType), 0x08,
				0xa6, 0x10,
				0xde, 0x1a,
				0x87, 0x09,
				0x23, 0x02,
				0x70, 0x05,
				0x6f, 0x05,
				0x2b, 0x03,
				0x9e, 0x04,
			},
			want: []int16{4262, 6878, 2439, 547, 1392, 1391, 811, 1182},
		},
		{
			name: "id 43 with negative sentinels",
			id:   43,
			input: []byte{
				0x2b, byte(v18.ShortArrayType), 0x0e,
				0x5a, 0x0d,
				0x65, 0x0d,
				0x68, 0x0d,
				0x65, 0x0f,
				0x79, 0x0f,
				0xdb, 0x0f,
				0xff, 0xff,
				0xff, 0xff,
				0xff, 0xff,
				0xff, 0xff,
				0xff, 0xff,
				0xff, 0xff,
				0xf2, 0x10,
				0x8a, 0x10,
			},
			want: []int16{3418, 3429, 3432, 3941, 3961, 4059, -1, -1, -1, -1, -1, -1, 4338, 4234},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)

			var parser v18.Parameter
			var got v18.Parameter

			if err := parser.Parse(r, &got, nil); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if got.ID() != tt.id {
				t.Errorf("ID() = %d, want %d", got.ID(), tt.id)
			}

			if got.Kind != v18.ShortArrayType {
				t.Errorf("Kind = %d, want %d", got.Kind, v18.ShortArrayType)
			}

			if got.Num != uint64(len(tt.want)) {
				t.Errorf("Num = %d, want %d", got.Num, len(tt.want))
			}

			var values []int16
			for _, value := range got.Int16ArrayValue() {
				values = append(values, value)
			}

			if !slices.Equal(values, tt.want) {
				t.Errorf("Int16ArrayValue() = %v, want %v", values, tt.want)
			}

			if r.Cursor != len(tt.input) {
				t.Errorf("Cursor = %d, want %d", r.Cursor, len(tt.input))
			}
		})
	}
}

func TestParseInt32ArrayParametersAndAccessor(t *testing.T) {
	tests := []struct {
		name  string
		id    uint8
		input []byte
		want  []int32
	}{
		{
			name: "id 8 with min sentinels",
			id:   8,
			input: []byte{
				0x08, byte(v18.CompressedIntArrayType), 0x03,
				0xde, 0xbb, 0xb3, 0xc3, 0x01,
				0xff, 0xff, 0xff, 0xff, 0x0f,
				0xff, 0xff, 0xff, 0xff, 0x0f,
			},
			want: []int32{204893935, -2147483648, -2147483648},
		},
		{
			name: "id 8 repeated minute durations",
			id:   8,
			input: []byte{
				0x08, byte(v18.CompressedIntArrayType), 0x08,
				0xc0, 0xa9, 0x07,
				0xc0, 0xa9, 0x07,
				0xc0, 0xa9, 0x07,
				0xc0, 0xa9, 0x07,
				0xc0, 0xa9, 0x07,
				0xc0, 0xa9, 0x07,
				0xc0, 0xa9, 0x07,
				0xe4, 0xb8, 0xb3, 0xc3, 0x01,
			},
			want: []int32{60000, 60000, 60000, 60000, 60000, 60000, 60000, 204893746},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)

			var parser v18.Parameter
			var got v18.Parameter

			if err := parser.Parse(r, &got, nil); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if got.ID() != tt.id {
				t.Errorf("ID() = %d, want %d", got.ID(), tt.id)
			}

			if got.Kind != v18.CompressedIntArrayType {
				t.Errorf("Kind = %d, want %d", got.Kind, v18.CompressedIntArrayType)
			}

			if got.Num != uint64(len(tt.want)) {
				t.Errorf("Num = %d, want %d", got.Num, len(tt.want))
			}

			var values []int32
			for _, value := range got.Int32ArrayValue() {
				values = append(values, value)
			}

			if !slices.Equal(values, tt.want) {
				t.Errorf("Int32ArrayValue() = %v, want %v", values, tt.want)
			}

			if r.Cursor != len(tt.input) {
				t.Errorf("Cursor = %d, want %d", r.Cursor, len(tt.input))
			}
		})
	}
}

func TestParseInt64ArrayParametersAndAccessor(t *testing.T) {
	tests := []struct {
		name  string
		id    uint8
		input []byte
		want  []int64
	}{
		{
			name: "id 4 timestamps",
			id:   4,
			input: []byte{
				0x04, byte(v18.CompressedLongArrayType), 0x08,
				0xf2, 0xd0, 0xd7, 0x92, 0xe4, 0xb4, 0xd0, 0xde, 0x11,
				0xce, 0xba, 0xd3, 0xf1, 0xe5, 0x96, 0xc8, 0xde, 0x11,
				0xf0, 0xc5, 0x88, 0xb9, 0x83, 0xb4, 0xd0, 0xde, 0x11,
				0x94, 0xa6, 0x97, 0xc4, 0xce, 0xa3, 0xcd, 0xde, 0x11,
				0xba, 0x86, 0xb3, 0xdb, 0x8a, 0xb5, 0xd0, 0xde, 0x11,
				0xa0, 0x8d, 0xe2, 0xe6, 0x8a, 0xb5, 0xd0, 0xde, 0x11,
				0xd8, 0xb0, 0x8a, 0xa6, 0x8b, 0xb5, 0xd0, 0xde, 0x11,
				0xb2, 0xba, 0xd0, 0xae, 0x8b, 0xb5, 0xd0, 0xde, 0x11,
			},
			want: []int64{
				639125025788195897,
				639106918439874215,
				639125012809322872,
				639118133759764874,
				639125030964715933,
				639125030976635728,
				639125031043025964,
				639125031051988633,
			},
		},
		{
			name: "id 4 duplicate timestamps",
			id:   4,
			input: []byte{
				0x04, byte(v18.CompressedLongArrayType), 0x05,
				0xd4, 0xac, 0x92, 0xbe, 0xf9, 0xb4, 0xd0, 0xde, 0x11,
				0x8a, 0x93, 0xfc, 0x98, 0xf7, 0xb4, 0xd0, 0xde, 0x11,
				0x8a, 0x93, 0xfc, 0x98, 0xf7, 0xb4, 0xd0, 0xde, 0x11,
				0xca, 0xfa, 0xc3, 0xd4, 0xda, 0xa0, 0xc4, 0xde, 0x11,
				0xe4, 0xad, 0xce, 0xb8, 0x8b, 0xb5, 0xd0, 0xde, 0x11,
			},
			want: []int64{
				639125028652337962,
				639125028344923333,
				639125028344923333,
				639098292638613157,
				639125031062457202,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)

			var parser v18.Parameter
			var got v18.Parameter

			if err := parser.Parse(r, &got, nil); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if got.ID() != tt.id {
				t.Errorf("ID() = %d, want %d", got.ID(), tt.id)
			}

			if got.Kind != v18.CompressedLongArrayType {
				t.Errorf("Kind = %d, want %d", got.Kind, v18.CompressedLongArrayType)
			}

			if got.Num != uint64(len(tt.want)) {
				t.Errorf("Num = %d, want %d", got.Num, len(tt.want))
			}

			var values []int64
			for _, value := range got.Int64ArrayValue() {
				values = append(values, value)
			}

			if !slices.Equal(values, tt.want) {
				t.Errorf("Int64ArrayValue() = %v, want %v", values, tt.want)
			}

			if r.Cursor != len(tt.input) {
				t.Errorf("Cursor = %d, want %d", r.Cursor, len(tt.input))
			}
		})
	}
}

func TestParseBooleanArrayParametersAndAccessor(t *testing.T) {

	tests := []struct {
		name  string
		id    uint8
		input []byte
		want  []bool
		num   uint64
	}{
		{
			name: "id 4 sparse packed flags",
			id:   4,
			input: []byte{
				0x04, byte(v18.BooleanArrayType), 0x86, 0x02,
				0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80,
				0x00, 0x00, 0x02, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			want: []bool{
				false, false, false, false, true, false, false, false, false,
				false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, true,
				false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, true,
				false, false, false, false, false, false, false, false, false,
				false, true, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false, false,
			},
			num: 262,
		},
		{
			name: "id 2 short packed flags",
			id:   2,
			input: []byte{
				0x02, byte(v18.BooleanArrayType), 0x09,
				0x37, 0x00,
			},
			want: []bool{true, true, true, false, true, true, false, false, false},
			num:  9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)

			var parser v18.Parameter
			var got v18.Parameter

			if err := parser.Parse(r, &got, nil); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if got.ID() != tt.id {
				t.Errorf("ID() = %d, want %d", got.ID(), tt.id)
			}

			if got.Kind != v18.BooleanArrayType {
				t.Errorf("Kind = %d, want %d", got.Kind, v18.BooleanArrayType)
			}

			if got.Num != tt.num {
				t.Errorf("Num = %d, want %d", got.Num, tt.num)
			}

			var values []bool
			for _, value := range got.BooleanArrayValue() {
				values = append(values, value)
			}

			if !slices.Equal(values, tt.want) {
				t.Errorf("BooleanArrayValue() = %v, want %v", values, tt.want)
			}
		})
	}
}

func TestParseByteArrayParametersAndAccessor(t *testing.T) {
	tests := []struct {
		name  string
		id    uint8
		input []byte
		want  []byte
	}{
		{
			name: "id 1",
			id:   1,
			input: []byte{
				0x01, byte(v18.ByteArrayType), 0x1e,
				0x03, 0xfb, 0x5d, 0x90, 0x5b, 0xd4, 0xa0, 0xde,
				0x08, 0x1f, 0x41, 0xaf, 0x3c, 0x63, 0x46, 0x3c,
				0x92, 0x33, 0x00, 0x00, 0xb0, 0x40, 0x8a, 0xda,
				0xb3, 0x3c, 0xe8, 0x46, 0x3e, 0x92,
			},
			want: []byte{
				3, 251, 93, 144, 91, 212, 160, 222,
				8, 31, 65, 175, 60, 99, 70, 60,
				146, 51, 0, 0, 176, 64, 138, 218,
				179, 60, 232, 70, 62, 146,
			},
		},
		{
			name: "id 7",
			id:   7,
			input: []byte{
				0x07, byte(v18.ByteArrayType), 0x06,
				0x4b, 0x12, 0x24, 0x49, 0x92, 0x00,
			},
			want: []byte{75, 18, 36, 73, 146, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)

			var parser v18.Parameter
			var got v18.Parameter

			if err := parser.Parse(r, &got, nil); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if got.ID() != tt.id {
				t.Errorf("ID() = %d, want %d", got.ID(), tt.id)
			}

			if got.Kind != v18.ByteArrayType {
				t.Errorf("Kind = %d, want %d", got.Kind, v18.ByteArrayType)
			}

			if got.Num != uint64(len(tt.want)) {
				t.Errorf("Num = %d, want %d", got.Num, len(tt.want))
			}

			var values []byte
			for _, value := range got.ByteArrayValue() {
				values = append(values, value)
			}

			if !slices.Equal(values, tt.want) {
				t.Errorf("ByteArrayValue() = %v, want %v", values, tt.want)
			}

			if r.Cursor != len(tt.input) {
				t.Errorf("Cursor = %d, want %d", r.Cursor, len(tt.input))
			}
		})
	}
}

func TestParseArrayParametersAndAccessor(t *testing.T) {
	tests := []struct {
		name  string
		id    uint8
		input []byte
		want  []any
	}{
		{
			name: "id 1",
			id:   1,
			input: []byte{
				0x01, byte(v18.ArrayType), 0x02,
				0x49, 0x03, 0x0e, 0x12, 0x16, 0x49, 0x05, 0x02,
				0x0e, 0x10, 0x12, 0x16,
			},
			want: []any{
				[]int32{7, 9, 11},
				[]int32{1, 7, 8, 9, 11},
			},
		},
		{
			name: "id 2",
			id:   2,
			input: []byte{
				0x02, byte(v18.ArrayType), 0x02,
				0x49, 0x08, 0x2e, 0x36, 0x3e, 0x38,
				0x4a, 0x42, 0x14, 0x00, 0x49, 0x03,
				0x0a, 0x38, 0x42, 0x49, 0x03, 0x06,
				0x38, 0x42, 0x49, 0x09, 0x46, 0x48,
				0x4a, 0x4c, 0x4e, 0x50, 0x52, 0x54,
				0x56, 0x49, 0x03, 0x04, 0x38, 0x42,
				0x49, 0x06, 0x34, 0x30, 0x32, 0x38,
				0x42, 0x4a, 0x49, 0x01, 0x30, 0x49,
				0x02, 0x3a, 0x3c, 0x49, 0x01, 0x34,
			},
			want: []any{
				[]int32{23, 27, 31, 28, 37, 33, 10, 0},
				[]int32{5, 28, 33},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)

			var parser v18.Parameter
			var got v18.Parameter

			if err := parser.Parse(r, &got, nil); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if got.ID() != tt.id {
				t.Errorf("ID() = %d, want %d", got.ID(), tt.id)
			}

			if got.Kind != v18.ArrayType {
				t.Errorf("Kind = %d, want %d", got.Kind, v18.ArrayType)
			}

			if got.Num != uint64(len(tt.want)) {
				t.Errorf("Num = %d, want %d", got.Num, len(tt.want))
			}

			var values []any
			for _, value := range got.ArrayValue() {
				values = append(values, value)
			}

			if !reflect.DeepEqual(values, tt.want) {
				t.Errorf("ArrayValue() = %v, want %v", values, tt.want)
			}
		})
	}
}
