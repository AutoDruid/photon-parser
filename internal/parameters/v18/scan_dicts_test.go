package v18_test

import (
	v18 "AutoDruid/photon-parser/internal/parameters/v18"
	"AutoDruid/photon-parser/internal/reader"
	"maps"
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
			name: "id 1",
			id:   1,
			input: []byte{
				0x01, byte(v18.DictionaryType), 0x09, 0x0a, 0x03,
				0x00, 0xc8, 0xba, 0xd3, 0xaf, 0x06, 0x02, 0x00, 0x04, 0x00,
			},
			want: map[any]any{
				int64(0): int64(855273124),
				int64(1): int64(0),
				int64(2): int64(0),
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

			if got.Kind != v18.DictionaryType {
				t.Errorf("Kind = %d, want %d", got.Kind, v18.DictionaryType)
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
