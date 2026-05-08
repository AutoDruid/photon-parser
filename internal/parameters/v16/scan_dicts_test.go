package v16_test

import (
	"maps"
	v16 "michelprogram/photon-parser/internal/parameters/v16"
	"michelprogram/photon-parser/internal/reader"
	"reflect"
	"testing"
)

func TestParseDictsParametersAndAccessor(t *testing.T) {
	tests := []struct {
		name  string
		id    uint8
		input []byte
		want  map[any]any
	}{
		{
			name: "int32 to int32 dictionary",
			input: []byte{
				0x00, 0x44,
				byte(v16.Int32Type), byte(v16.Int32Type), // keyType, valueType
				0x00, 0x02, // size = 2
				0x00, 0x00, 0x00, 0x01, // key = 1
				0x00, 0x00, 0x28, 0xC1, // value = 10433
				0x00, 0x00, 0x00, 0x02, // key = 2
				0x00, 0x00, 0x02, 0x21, // value = 545
			},
			want: map[any]any{
				int64(1): int64(10433),
				int64(2): int64(545),
			},
		},
		{
			name: "empty dictionary",
			input: []byte{
				0x00, 0x44,
				byte(v16.Int32Type), byte(v16.Int32Type),
				0x00, 0x00, // size = 0
			},
			want: map[any]any{},
		},
		{
			name: "string to boolean dictionary",
			input: []byte{
				0x00, 0x44,
				byte(v16.StringType), byte(v16.BooleanType),
				0x00, 0x02, // size = 2
				0x00, 0x01, 'a', // key = "a"
				0x01,                 // value = true
				0x00, 0x02, 'b', 'b', // key = "bb"
				0x00, // value = false
			},
			want: map[any]any{
				"a":  true,
				"bb": false,
			},
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

			if got.ID() != tt.id {
				t.Errorf("ID() = %d, want %d", got.ID(), tt.id)
			}

			if got.Kind != v16.DictionaryType {
				t.Errorf("Kind = %d, want %d", got.Kind, v16.DictionaryType)
			}

			if got.Num != uint64(len(tt.want)) {
				t.Errorf("Num = %d, want %d", got.Num, len(tt.want))
			}

			values := make(map[any]any, got.Num)
			maps.Insert(values, got.DictionaryValue())

			if !reflect.DeepEqual(values, tt.want) {
				t.Errorf("DictionaryValue() = %v, want %v", values, tt.want)
			}
		})
	}
}
